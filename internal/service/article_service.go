package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/redis"
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// generateSlug 根据标题生成 URL slug（仅保留 ASCII 字母数字，避免中文导致 URL 乱码）
func generateSlug(title string) string {
	// 转换为小写
	slug := strings.ToLower(title)
	// 只保留 ASCII 小写字母、数字、空格和连字符；中文等非 ASCII 字符全部移除
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, slug)
	// 替换空格为连字符
	slug = strings.TrimSpace(slug)
	re := regexp.MustCompile(`\s+`)
	slug = re.ReplaceAllString(slug, "-")
	// 去除连续连字符
	re2 := regexp.MustCompile(`-+`)
	slug = re2.ReplaceAllString(slug, "-")
	// 去除首尾连字符
	slug = strings.Trim(slug, "-")
	// 如果 slug 为空（纯中文标题），使用时间戳生成
	if slug == "" {
		slug = fmt.Sprintf("post-%d", time.Now().Unix())
	}
	// 限制长度
	if len(slug) > 200 {
		slug = slug[:200]
	}
	return slug
}

// calculateReadingTime 计算阅读时间（分钟）
// 中文约 400 字/分钟，英文约 200 词/分钟
func calculateReadingTime(content string) int {
	if content == "" {
		return 1
	}
	// 去除 markdown 标记
	cleaned := content
	replacements := []string{"#", "*", "`", ">", "[", "]", "(", ")", "!", "-"}
	for _, r := range replacements {
		cleaned = strings.ReplaceAll(cleaned, r, "")
	}
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return 1
	}

	// 计算中文字符数
	runeCount := len([]rune(cleaned))
	// 计算英文单词数
	wordCount := len(strings.Fields(cleaned))

	// 混合计算：中文字符按 400 字/分钟，英文单词按 200 词/分钟
	minutes := (runeCount / 400) + (wordCount / 200)
	if minutes < 1 {
		return 1
	}
	return minutes
}

// articleService 文章服务实现
type articleService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
	visitRepo    repository.VisitRepository
	redisClient  *redis.Client
	scheduler    ArticlePublishScheduler
}

// NewArticleService 创建文章服务
func NewArticleService(
	articleRepo repository.ArticleRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
	visitRepo repository.VisitRepository,
	redisClient *redis.Client,
) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		visitRepo:    visitRepo,
		redisClient:  redisClient,
	}
}

func (s *articleService) SetPublishScheduler(scheduler ArticlePublishScheduler) {
	s.scheduler = scheduler
}

func (s *articleService) HandleScheduledArticlePublished(article *entity.Article) {
	if article != nil {
		s.invalidateArticleCache(article.Slug)
		return
	}
	s.invalidateAllArticleListCache()
}

// GetArticleList 获取文章列表（前台）
func (s *articleService) GetArticleList(req *request.ArticleListRequest) (*response.PageResponse, error) {
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("article:list:%d:%d:%d:%d:%s",
			req.Page, req.GetPageSize(), req.Category, req.Tag, req.Keyword)
		var cachedPage response.PageResponse
		if err := s.redisClient.GetJSON(context.Background(), cacheKey, &cachedPage); err == nil {
			return &cachedPage, nil
		}
	}

	list, total, err := s.articleRepo.ListPublished(req.GetOffset(), req.GetPageSize(), req.Category, req.Tag, req.Keyword)
	if err != nil {
		return nil, fmt.Errorf("获取文章列表失败, %w", err)
	}

	pageResp := response.NewPageResponse(list, total, req.Page, req.GetPageSize())

	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("article:list:%d:%d:%d:%d:%s",
			req.Page, req.GetPageSize(), req.Category, req.Tag, req.Keyword)
		go func() {
			defer func() { _ = recover() }()
			_ = s.redisClient.SetJSON(context.Background(), cacheKey, pageResp, 5*time.Minute)
		}()
	}

	return pageResp, nil
}

// GetArticleDetail 获取文章详情（前台）
func (s *articleService) GetArticleDetail(slug string, clientIP string) (*response.ArticleDetailResponse, error) {
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("article:detail:%s", slug)
		var cachedDetail response.ArticleDetailResponse
		if err := s.redisClient.GetJSON(context.Background(), cacheKey, &cachedDetail); err == nil {
			go s.incrementViewCountAsync(cachedDetail.ID, clientIP)
			return &cachedDetail, nil
		}
	}

	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("获取文章详情失败, %w", err)
	}
	if article == nil {
		return nil, bizerrors.New(bizerrors.CodeArticleNotFound, bizerrors.GetMessage(bizerrors.CodeArticleNotFound))
	}

	viewIncremented := false
	if clientIP != "" && s.visitRepo != nil {
		since := time.Now().Add(-24 * time.Hour)
		hasVisited, _ := s.visitRepo.HasVisited(article.ID, clientIP, since)
		if !hasVisited {
			_ = s.articleRepo.IncrementViewCount(article.ID)
			_ = s.visitRepo.Create(&entity.VisitLog{
				ArticleID: article.ID,
				IP:        clientIP,
			})
			viewIncremented = true
		}
	} else {
		_ = s.articleRepo.IncrementViewCount(article.ID)
		viewIncremented = true
	}

	viewCount := article.ViewCount
	if viewIncremented {
		viewCount++
	}

	cats := response.CategoryResponse{}
	if article.Category.ID > 0 {
		cats.ID = article.Category.ID
		cats.Name = article.Category.Name
		cats.Slug = article.Category.Slug
	}
	tags := make([]response.TagResponse, 0, len(article.Tags))
	for _, t := range article.Tags {
		tags = append(tags, response.TagResponse{ID: t.ID, Name: t.Name, Slug: t.Slug})
	}

	detail := &response.ArticleDetailResponse{
		ID:           article.ID,
		Title:        article.Title,
		Slug:         article.Slug,
		Content:      article.Content,
		Summary:      article.Summary,
		Cover:        article.Cover,
		Category:     cats,
		Tags:         tags,
		ViewCount:    viewCount,
		CommentCount: article.CommentCount,
		Status:       article.Status,
		IsTop:        article.IsTop,
		ReadingTime:  article.ReadingTime,
		SEOTitle:     article.SEOTitle,
		SEODesc:      article.SEODescription,
		SEOKeywords:  article.SEOKeywords,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}

	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("article:detail:%s", slug)
		go func() {
			defer func() { _ = recover() }()
			_ = s.redisClient.SetJSON(context.Background(), cacheKey, detail, 10*time.Minute)
		}()
	}

	return detail, nil
}

func (s *articleService) invalidateArticleCache(slug string) {
	if s.redisClient == nil {
		return
	}
	ctx := context.Background()
	var keys []string
	if slug != "" {
		keys = append(keys, fmt.Sprintf("article:detail:%s", slug))
	}
	keys = append(keys,
		"category:list",
		"tag:list",
		"article:archives",
	)
	go func() {
		defer func() { _ = recover() }()
		_ = s.redisClient.Del(ctx, keys...)
	}()
	// 文章内容变化也会影响列表摘要/搜索展示，列表缓存一并失效
	go func() {
		defer func() { _ = recover() }()
		_ = s.redisClient.DelByPattern(ctx, "article:list:*")
	}()
}

func (s *articleService) invalidateAllArticleListCache() {
	if s.redisClient == nil {
		return
	}
	ctx := context.Background()
	keys := []string{
		"category:list",
		"tag:list",
		"article:archives",
	}
	go func() {
		defer func() { _ = recover() }()
		_ = s.redisClient.Del(ctx, keys...)
	}()
	// 清除文章列表分页缓存（article:list:*）
	go func() {
		defer func() { _ = recover() }()
		_ = s.redisClient.DelByPattern(ctx, "article:list:*")
	}()
}

func (s *articleService) refreshArticleSchedule(article *entity.Article) {
	if s.scheduler == nil || article == nil {
		return
	}
	if err := s.scheduler.RefreshArticle(article); err != nil {
		logger.Warnf("刷新文章定时任务失败, id: %d, error: %v", article.ID, err)
	}
}

func (s *articleService) unscheduleArticle(id uint) {
	if s.scheduler == nil {
		return
	}
	s.scheduler.UnscheduleArticle(id)
}

func (s *articleService) incrementViewCountAsync(articleID uint, clientIP string) {
	if clientIP != "" && s.visitRepo != nil {
		since := time.Now().Add(-24 * time.Hour)
		hasVisited, _ := s.visitRepo.HasVisited(articleID, clientIP, since)
		if !hasVisited {
			_ = s.articleRepo.IncrementViewCount(articleID)
			_ = s.visitRepo.Create(&entity.VisitLog{
				ArticleID: articleID,
				IP:        clientIP,
			})
		}
	} else {
		_ = s.articleRepo.IncrementViewCount(articleID)
	}
}

// GetArticleArchives 获取文章归档
func (s *articleService) GetArticleArchives() ([]response.ArchiveResponse, error) {
	// 获取所有已发布文章，按创建时间降序排序（大上限避免归档遗漏）
	articles, _, err := s.articleRepo.ListPublished(0, 100000, 0, 0, "")
	if err != nil {
		return nil, fmt.Errorf("获取文章归档失败, %w", err)
	}

	// 按年份分组
	yearMap := make(map[int][]response.ArchiveArticleResponse)
	for _, article := range articles {
		year := article.CreatedAt.Year()
		yearMap[year] = append(yearMap[year], response.ArchiveArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Slug:      article.Slug,
			CreatedAt: article.CreatedAt,
		})
	}

	// 转换为响应格式，并按年份降序排序
	var result []response.ArchiveResponse
	for year, arts := range yearMap {
		// 确保每个年份内的文章按创建时间降序排序
		sort.Slice(arts, func(i, j int) bool {
			return arts[i].CreatedAt.After(arts[j].CreatedAt)
		})
		result = append(result, response.ArchiveResponse{
			Year:     year,
			Articles: arts,
		})
	}

	// 按年份降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Year > result[j].Year
	})

	return result, nil
}

// GetAdminArticleList 获取文章列表（后台）
func (s *articleService) GetAdminArticleList(req *request.ArticleListRequest) (*response.PageResponse, error) {
	if req.Status != "" && req.Status != entity.ArticleStatusPublished && req.Status != entity.ArticleStatusDraft && req.Status != entity.ArticleStatusScheduled {
		return nil, bizerrors.New(bizerrors.CodeInvalidParams, "无效的文章状态")
	}
	list, total, err := s.articleRepo.ListAll(req.GetOffset(), req.GetPageSize(), req.Status, req.Keyword, req.Category)
	if err != nil {
		return nil, fmt.Errorf("获取后台文章列表失败, %w", err)
	}
	return response.NewPageResponse(list, total, req.Page, req.GetPageSize()), nil
}

// GetAdminArticleDetail 获取文章详情（后台）
func (s *articleService) GetAdminArticleDetail(id uint) (*response.AdminArticleDetailResponse, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取后台文章详情失败, %w", err)
	}
	if article == nil {
		return nil, bizerrors.New(bizerrors.CodeArticleNotFound, bizerrors.GetMessage(bizerrors.CodeArticleNotFound))
	}

	// 从 Tags 中提取 tag IDs
	var tagIDs []uint
	for _, tag := range article.Tags {
		tagIDs = append(tagIDs, tag.ID)
	}

	return &response.AdminArticleDetailResponse{
		ID:             article.ID,
		Title:          article.Title,
		Slug:           article.Slug,
		Content:        article.Content,
		Summary:        article.Summary,
		Cover:          article.Cover,
		Status:         article.Status,
		IsTop:          article.IsTop,
		CategoryID:     article.CategoryID,
		TagIDs:         tagIDs,
		ScheduledAt:    article.ScheduledAt,
		SEOTitle:       article.SEOTitle,
		SEODescription: article.SEODescription,
		SEOKeywords:    article.SEOKeywords,
	}, nil
}

// CreateArticle 创建文章
func (s *articleService) CreateArticle(req *request.CreateArticleRequest) (uint, error) {
	// 查询标签
	var tags []entity.Tag
	if len(req.TagIDs) > 0 {
		tagList, err := s.tagRepo.FindByIDs(req.TagIDs)
		if err != nil {
			return 0, fmt.Errorf("查询标签失败, %w", err)
		}
		for _, t := range tagList {
			tags = append(tags, *t)
		}
	}

	// 生成或使用 Slug
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	article := &entity.Article{
		Title:          req.Title,
		Slug:           slug,
		Content:        req.Content,
		Summary:        req.Summary,
		Cover:          req.Cover,
		CategoryID:     req.CategoryID,
		Status:         req.Status,
		IsTop:          req.IsTop,
		Tags:           tags,
		ReadingTime:    calculateReadingTime(req.Content),
		ScheduledAt:    req.ScheduledAt,
		SEOTitle:       req.SEOTitle,
		SEODescription: req.SEODescription,
		SEOKeywords:    req.SEOKeywords,
	}

	if article.Status == "" {
		article.Status = entity.ArticleStatusDraft
	}

	// 迁移历史文章时保留原始发布时间
	if req.CreatedAt != nil {
		article.CreatedAt = *req.CreatedAt
	}

	err := s.articleRepo.Create(article)
	if err != nil {
		return 0, fmt.Errorf("创建文章失败, %w", err)
	}

	logger.Infof("文章创建成功, id: %d, title: %s, tags: %d, slug: %s", article.ID, article.Title, len(tags), article.Slug)

	s.invalidateAllArticleListCache()
	s.refreshArticleSchedule(article)

	return article.ID, nil
}

// UpdateArticle 更新文章
func (s *articleService) UpdateArticle(id uint, req *request.UpdateArticleRequest) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询文章失败, %w", err)
	}
	if article == nil {
		return bizerrors.New(bizerrors.CodeArticleNotFound, bizerrors.GetMessage(bizerrors.CodeArticleNotFound))
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
		article.ReadingTime = calculateReadingTime(req.Content)
	}
	if req.Summary != nil {
		article.Summary = *req.Summary
	}
	if req.Cover != nil {
		article.Cover = *req.Cover
	}
	if req.CategoryID != 0 {
		article.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		article.Status = req.Status
	}
	article.IsTop = req.IsTop

	// Slug：支持自定义，为空时自动根据标题生成
	if req.Slug != "" {
		article.Slug = req.Slug
	} else if req.Title != "" {
		article.Slug = generateSlug(req.Title)
	}

	// 定时发布
	article.ScheduledAt = req.ScheduledAt

	// SEO 字段
	if req.SEOTitle != nil {
		article.SEOTitle = *req.SEOTitle
	}
	if req.SEODescription != nil {
		article.SEODescription = *req.SEODescription
	}
	if req.SEOKeywords != nil {
		article.SEOKeywords = *req.SEOKeywords
	}

	// 更新标签关联
	if req.TagIDs != nil {
		var tags []entity.Tag
		if len(req.TagIDs) > 0 {
			tagList, err := s.tagRepo.FindByIDs(req.TagIDs)
			if err != nil {
				return fmt.Errorf("查询标签失败, %w", err)
			}
			for _, t := range tagList {
				tags = append(tags, *t)
			}
		}
		// 清除旧关联，建立新关联
		if err := s.articleRepo.UpdateTags(article, tags); err != nil {
			return fmt.Errorf("更新文章标签失败, %w", err)
		}
	}

	if err := s.articleRepo.Update(article); err != nil {
		return fmt.Errorf("更新文章失败, %w", err)
	}

	logger.Infof("文章更新成功, id: %d", id)

	s.invalidateArticleCache(article.Slug)
	s.refreshArticleSchedule(article)

	return nil
}

// DeleteArticle 删除文章
func (s *articleService) DeleteArticle(id uint) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询文章失败, %w", err)
	}
	if article == nil {
		return bizerrors.New(bizerrors.CodeArticleNotFound, bizerrors.GetMessage(bizerrors.CodeArticleNotFound))
	}

	if err := s.articleRepo.Delete(id); err != nil {
		return fmt.Errorf("删除文章失败, %w", err)
	}

	logger.Infof("文章删除成功, id: %d", id)

	s.invalidateArticleCache(article.Slug)
	s.unscheduleArticle(id)

	return nil
}

// BatchDeleteArticles 批量删除文章
func (s *articleService) BatchDeleteArticles(ids []uint) error {
	if err := s.articleRepo.BatchDelete(ids); err != nil {
		return fmt.Errorf("批量删除文章失败, %w", err)
	}

	logger.Infof("批量删除文章成功, ids: %v", ids)

	s.invalidateAllArticleListCache()
	for _, id := range ids {
		s.unscheduleArticle(id)
	}

	return nil
}

// regexpCache 预编译正则缓存（避免 MustCompile panic，启动时安全编译）
var (
	regexpOnce    sync.Once
	reFenced      *regexp.Regexp
	reInlineCode  *regexp.Regexp
	reHTMLComment *regexp.Regexp
	reHTMLTag     *regexp.Regexp
	reImage       *regexp.Regexp
	reLink        *regexp.Regexp
	reLinePrefix  *regexp.Regexp
	// 下列正则改为「从外向内逐条剥离」的独立写法，完全避免反向引用（RE2 不支持 \1）
	reTripleStar  *regexp.Regexp
	reDoubleStar  *regexp.Regexp
	reSingleStar  *regexp.Regexp
	reTripleTilde *regexp.Regexp
	reDoubleTilde *regexp.Regexp
	reTripleUnder *regexp.Regexp
	reDoubleUnder *regexp.Regexp
	reSingleUnder *regexp.Regexp
	reWS          *regexp.Regexp
)

func compileRegexps() {
	regexpOnce.Do(func() {
		reFenced = mustOrEmpty("(?s)```" + `[^\n]*\n?(.*?)` + "```")
		reInlineCode = mustOrEmpty("`([^`]*)`")
		reHTMLComment = mustOrEmpty(`(?s)<!--.*?-->`)
		reHTMLTag = mustOrEmpty(`<[^>]+>`)
		reImage = mustOrEmpty(`!\[([^\]]*)\]\([^)]*\)`)
		reLink = mustOrEmpty(`\[([^\]]*)\]\([^)]*\)`)
		reLinePrefix = mustOrEmpty("(?m)^\\s{0,3}(#{1,6}\\s*|>\\s*|[-*+]\\s+|\\d+\\.\\s+)")
		reTripleStar = mustOrEmpty(`\*\*\*([^*]+?)\*\*\*`)
		reDoubleStar = mustOrEmpty(`\*\*([^*]+?)\*\*`)
		reSingleStar = mustOrEmpty(`(^|[^*])\*([^*\n]+?)\*([^*]|$)`)
		reTripleTilde = mustOrEmpty(`~~~([^~]+?)~~~`)
		reDoubleTilde = mustOrEmpty(`~~([^~]+?)~~`)
		reTripleUnder = mustOrEmpty(`___([^_]+?)___`)
		reDoubleUnder = mustOrEmpty(`__([^_]+?)__`)
		reSingleUnder = mustOrEmpty(`(^|[^_\w])_([^_\n]+?)_([^_\w]|$)`)
		reWS = mustOrEmpty(`\s+`)
	})
}

// mustOrEmpty 编译 regexp；失败时返回一个永不匹配的空壳，保证服务不 panic
func mustOrEmpty(expr string) *regexp.Regexp {
	r, err := regexp.Compile(expr)
	if err != nil {
		// 构造一个不可能匹配的表达式（空字符串会被 RE2 视为错误）
		return regexp.MustCompile(`a^`)
	}
	return r
}

// stripMarkdown 粗粒度清除 Markdown 标记，保留纯文本（用于搜索上下文片段）
func stripMarkdown(src string) string {
	compileRegexps()
	if src == "" {
		return ""
	}
	s := src
	s = reFenced.ReplaceAllString(s, " $1 ")
	s = reInlineCode.ReplaceAllString(s, " $1 ")
	s = reHTMLComment.ReplaceAllString(s, "")
	s = reHTMLTag.ReplaceAllString(s, "")
	s = reImage.ReplaceAllString(s, "$1")
	s = reLink.ReplaceAllString(s, "$1")
	s = reLinePrefix.ReplaceAllString(s, "")
	// 嵌套的 ***...*** / **...** / *...* 从多到少剥离
	for reTripleStar.MatchString(s) {
		s = reTripleStar.ReplaceAllString(s, "$1")
	}
	for reDoubleStar.MatchString(s) {
		s = reDoubleStar.ReplaceAllString(s, "$1")
	}
	for reSingleStar.MatchString(s) {
		s = reSingleStar.ReplaceAllString(s, "$1$2$3")
	}
	for reTripleTilde.MatchString(s) {
		s = reTripleTilde.ReplaceAllString(s, "$1")
	}
	for reDoubleTilde.MatchString(s) {
		s = reDoubleTilde.ReplaceAllString(s, "$1")
	}
	for reTripleUnder.MatchString(s) {
		s = reTripleUnder.ReplaceAllString(s, "$1")
	}
	for reDoubleUnder.MatchString(s) {
		s = reDoubleUnder.ReplaceAllString(s, "$1")
	}
	for reSingleUnder.MatchString(s) {
		s = reSingleUnder.ReplaceAllString(s, "$1$2$3")
	}
	// 最后兜底：把孤立的 * ~ _ 符号去掉（避免残留的标记）
	s = strings.Map(func(r rune) rune {
		switch r {
		case '*', '~':
			return -1
		}
		return r
	}, s)
	s = reWS.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// keywordPattern 按空格拆分关键词做 OR 匹配的不区分大小写正则（支持 Unicode）
func keywordPattern(keyword string) (*regexp.Regexp, error) {
	keywords := strings.Fields(keyword)
	if len(keywords) == 0 {
		return nil, nil
	}
	parts := make([]string, 0, len(keywords))
	for _, k := range keywords {
		if k == "" {
			continue
		}
		parts = append(parts, regexp.QuoteMeta(k))
	}
	if len(parts) == 0 {
		return nil, nil
	}
	return regexp.Compile("(?i)(" + strings.Join(parts, "|") + ")")
}

// buildSearchSnippet 从文章正文中抽取包含关键词的上下文片段，前后用 "[...]" 省略。
// radius 为关键词前后保留的 rune 数（中文近似为汉字数）。
// 若文章 Summary 本身命中关键词，返回 ""（让前端直接用 Summary，不再截断）。
func buildSearchSnippet(content, summary, keyword string, radius int) string {
	if keyword == "" {
		return ""
	}
	kwLower := strings.ToLower(keyword)
	sumLower := strings.ToLower(summary)
	for _, w := range strings.Fields(kwLower) {
		if w != "" && strings.Contains(sumLower, w) {
			return ""
		}
	}
	clean := stripMarkdown(content)
	if clean == "" {
		return ""
	}
	re, err := keywordPattern(keyword)
	if err != nil || re == nil {
		return ""
	}
	loc := re.FindStringIndex(clean)
	if loc == nil {
		return ""
	}
	runes := []rune(clean)
	// 把 byte 范围映射到 rune 范围（精确扫描）
	hitStart, hitEnd := -1, -1
	cur := 0
	for i := 0; i < len(runes); i++ {
		next := cur + len(string(runes[i:i+1]))
		if cur <= loc[0] && loc[0] < next && hitStart == -1 {
			hitStart = i
		}
		if cur <= loc[1] && loc[1] <= next {
			hitEnd = i + 1
		}
		if hitStart >= 0 && hitEnd >= 0 {
			break
		}
		cur = next
	}
	if hitStart < 0 {
		hitStart = 0
	}
	if hitEnd < 0 {
		hitEnd = len(runes)
	}
	left := hitStart - radius
	if left < 0 {
		left = 0
	}
	right := hitEnd + radius
	if right > len(runes) {
		right = len(runes)
	}
	segment := strings.TrimSpace(string(runes[left:right]))
	if segment == "" {
		return ""
	}
	var b strings.Builder
	if left > 0 {
		b.WriteString("[...]")
	}
	b.WriteString(segment)
	if right < len(runes) {
		b.WriteString("[...]")
	}
	return b.String()
}

// Search 搜索文章（标题、内容、摘要模糊搜索），返回带 SearchSnippet 的 SearchArticleResponse 列表
func (s *articleService) Search(keyword string, page, pageSize int) (*response.PageResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	articles, total, err := s.articleRepo.Search(keyword, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("搜索文章失败, %w", err)
	}

	const snippetRadius = 45 // 关键词前后约 45 个汉字
	items := make([]*response.SearchArticleResponse, 0, len(articles))
	for _, a := range articles {
		cats := response.CategoryResponse{}
		if a.Category.ID > 0 {
			cats.ID = a.Category.ID
			cats.Name = a.Category.Name
			cats.Slug = a.Category.Slug
		}
		tags := make([]response.TagResponse, 0, len(a.Tags))
		for _, t := range a.Tags {
			tags = append(tags, response.TagResponse{ID: t.ID, Name: t.Name, Slug: t.Slug})
		}
		items = append(items, &response.SearchArticleResponse{
			ID:            a.ID,
			Title:         a.Title,
			Slug:          a.Slug,
			Summary:       a.Summary,
			SearchSnippet: buildSearchSnippet(a.Content, a.Summary, keyword, snippetRadius),
			Cover:         a.Cover,
			Category:      cats,
			Tags:          tags,
			ViewCount:     a.ViewCount,
			CommentCount:  a.CommentCount,
			ReadingTime:   a.ReadingTime,
			CreatedAt:     a.CreatedAt,
		})
	}

	return response.NewPageResponse(items, total, page, pageSize), nil
}
