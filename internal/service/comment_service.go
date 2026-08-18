package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/config"
	"blog/pkg/email"
	bizerrors "blog/pkg/errors"
	"blog/pkg/gravatar"
	"blog/pkg/logger"
	"blog/pkg/ua"
	"fmt"
	"strconv"
	"strings"
)

// commentService 评论服务实现
type commentService struct {
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
	emailSvc    email.EmailService
	config      *config.Config
}

// NewCommentService 创建评论服务
func NewCommentService(commentRepo repository.CommentRepository, articleRepo repository.ArticleRepository, userRepo repository.UserRepository, emailSvc email.EmailService, config *config.Config) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
		emailSvc:    emailSvc,
		config:      config,
	}
}

// isBlogger 判断 userID 是否为博主（硬编码配置，不查用户表）
func (s *commentService) isBlogger(userID *uint) bool {
	if s.config == nil || userID == nil || *userID == 0 {
		return false
	}
	return *userID == s.config.Blogger.UserID
}

// convertToCommentResponse 将 entity.Comment 转换为 response.CommentResponse
// 博主身份判断：UserID 等于配置中博主虚拟 ID 即标记为博主
func (s *commentService) convertToCommentResponse(comment *entity.Comment) response.CommentResponse {
	resp := response.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		Nickname:  comment.Nickname,
		Email:     comment.Email,
		Website:   comment.Website,
		Avatar:    comment.Avatar,
		Status:    comment.Status,
		LikeCount: comment.LikeCount,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
		IsAdmin:   s.isBlogger(comment.UserID),
	}

	if comment.Article.ID > 0 {
		resp.Article = response.ArticleBriefResponse{
			ID:    comment.Article.ID,
			Title: comment.Article.Title,
			Slug:  comment.Article.Slug,
		}
	}

	if resp.Avatar == "" && resp.Email != "" {
		resp.Avatar = gravatar.GetAvatarURLByEmail(resp.Email, 80)
	}

	if comment.ReplyToNickname != "" {
		resp.ReplyTo = comment.ReplyToNickname
	}

	// 客户端信息优先直接取已经存好的列（前端精确检测/后端兜底），
	// 历史空值数据才回退到解析 User-Agent
	hasStoredInfo := (comment.OS != "" && comment.OS != "未知") || (comment.Browser != "" && comment.Browser != "未知")
	if hasStoredInfo {
		if comment.OS != "未知" {
			resp.OS = comment.OS
			resp.OSVersion = comment.OSVersion
		}
		if comment.Browser != "未知" {
			resp.Browser = comment.Browser
			resp.BrowserVersion = comment.BrowserVersion
		}
	} else if strings.TrimSpace(comment.UserAgent) != "" {
		uaInfo := ua.Parse(comment.UserAgent)
		if uaInfo.OS != "未知" {
			resp.OS = uaInfo.OS
			resp.OSVersion = uaInfo.OSVersion
		}
		if uaInfo.Browser != "未知" {
			resp.Browser = uaInfo.Browser
			resp.BrowserVersion = uaInfo.BrowserVersion
		}
	}

	return resp
}

// buildCommentTree 构建评论响应树并注入博主标识
func (s *commentService) buildCommentResponseList(rootComments []*entity.Comment) []response.CommentResponse {
	result := make([]response.CommentResponse, 0, len(rootComments))
	for _, root := range rootComments {
		rootResp := s.convertToCommentResponse(root)
		if len(root.Replies) > 0 {
			repliesResp := make([]response.CommentResponse, 0, len(root.Replies))
			for i := range root.Replies {
				reply := &root.Replies[i]
				replyResp := s.convertToCommentResponse(reply)
				repliesResp = append(repliesResp, replyResp)
			}
			rootResp.Replies = repliesResp
		}
		result = append(result, rootResp)
	}
	return result
}

// GetCommentsByArticle 获取文章评论列表（支持 slug 或数字 ID）
func (s *commentService) GetCommentsByArticle(articleParam string, req *request.CommentListRequest) (*response.PageResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var articleID uint
	if id, err := strconv.ParseUint(articleParam, 10, 32); err == nil {
		articleID = uint(id)
	} else {
		article, err := s.articleRepo.FindBySlug(articleParam)
		if err != nil {
			return nil, fmt.Errorf("根据slug查询文章失败, %w", err)
		}
		if article == nil {
			return nil, bizerrors.New(bizerrors.CodeArticleNotFound, bizerrors.GetMessage(bizerrors.CodeArticleNotFound))
		}
		articleID = article.ID
	}

	sortBy := strings.ToLower(strings.TrimSpace(req.SortBy))
	if sortBy != "asc" && sortBy != "desc" && sortBy != "hot" {
		sortBy = "desc"
	}

	list, total, err := s.commentRepo.ListByArticleID(articleID, (req.Page-1)*req.PageSize, req.PageSize, sortBy)
	if err != nil {
		return nil, fmt.Errorf("获取文章评论列表失败, %w", err)
	}

	respList := s.buildCommentResponseList(list)

	return response.NewPageResponse(respList, total, req.Page, req.PageSize), nil
}

// CreateComment 创建评论
// userID: 当前登录用户 ID，0 表示访客
// 博主身份判定条件：userID > 0（已通过登录按钮登录）
func (s *commentService) CreateComment(req *request.CreateCommentRequest, userID uint, ip, userAgent string) (uint, error) {
	if strings.TrimSpace(req.Content) == "" {
		return 0, bizerrors.New(bizerrors.CodeInvalidParams, "评论内容不能为空")
	}

	if strings.TrimSpace(req.Nickname) == "" {
		return 0, bizerrors.New(bizerrors.CodeInvalidParams, "昵称不能为空")
	}

	var article *entity.Article
	var err error
	if req.ArticleSlug != "" {
		article, err = s.articleRepo.FindBySlug(req.ArticleSlug)
	} else if req.ArticleID > 0 {
		article, err = s.articleRepo.FindByID(req.ArticleID)
	} else {
		return 0, bizerrors.New(bizerrors.CodeArticleNotFound, bizerrors.GetMessage(bizerrors.CodeArticleNotFound))
	}
	if err != nil {
		return 0, fmt.Errorf("查询文章失败, %w", err)
	}
	if article == nil {
		return 0, bizerrors.New(bizerrors.CodeArticleNotFound, bizerrors.GetMessage(bizerrors.CodeArticleNotFound))
	}

	comment := &entity.Comment{
		Content:   req.Content,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Website:   req.Website,
		ArticleID: article.ID,
		ParentID:  req.ParentID,
		ReplyToID: req.ReplyToID,
		Status:    "approved",
		IP:        ip,
		UserAgent: userAgent,
	}

	// 客户端信息：优先使用前端通过 JS 精确检测的结果（能区分 Win11/Win10），
	// 空值才用后端 UA 解析兜底
	uaInfo := ua.Parse(userAgent)
	switch {
	case strings.TrimSpace(req.OS) != "" && req.OS != "未知":
		comment.OS = strings.TrimSpace(req.OS)
		comment.OSVersion = strings.TrimSpace(req.OSVersion)
	default:
		if uaInfo.OS != "未知" {
			comment.OS = uaInfo.OS
			comment.OSVersion = uaInfo.OSVersion
		}
	}
	switch {
	case strings.TrimSpace(req.Browser) != "" && req.Browser != "未知":
		comment.Browser = strings.TrimSpace(req.Browser)
		comment.BrowserVersion = strings.TrimSpace(req.BrowserVersion)
	default:
		if uaInfo.Browser != "未知" {
			comment.Browser = uaInfo.Browser
			comment.BrowserVersion = uaInfo.BrowserVersion
		}
	}

	if userID > 0 {
		// 博主账号：强制覆盖前端传值，昵称/头像/邮箱统一从 user 表读取（与主页 /api/v1/user/info 同源），
		// 用户表无记录时再回退配置文件，防止冒充
		if s.config != nil && userID == s.config.Blogger.UserID {
			b := s.config.Blogger
			nickname := b.Nickname
			email := b.Email
			avatar := b.Avatar
			if u, dbErr := s.userRepo.FindByID(userID); dbErr == nil && u != nil {
				if u.Nickname != "" {
					nickname = u.Nickname
				}
				if u.Email != "" {
					email = u.Email
				}
				if u.Avatar != "" {
					avatar = u.Avatar
				}
			}
			comment.UserID = &userID
			comment.Nickname = nickname
			comment.Email = email
			comment.Avatar = avatar
		} else {
			// 其他登录用户：查用户表，空值才填充
			user, err := s.userRepo.FindByID(userID)
			if err == nil && user != nil {
				comment.UserID = &userID
				if comment.Nickname == "" {
					comment.Nickname = user.Nickname
				}
				if comment.Email == "" {
					comment.Email = user.Email
				}
				if comment.Avatar == "" {
					comment.Avatar = user.Avatar
				}
			}
		}
	}

	if comment.Avatar == "" && comment.Email != "" {
		comment.Avatar = gravatar.GetAvatarURLByEmail(comment.Email, 80)
	}

	err = s.commentRepo.Create(comment)
	if err != nil {
		return 0, fmt.Errorf("创建评论失败, %w", err)
	}

	// 同步文章评论数（新评论默认 approved）
	if err := s.articleRepo.UpdateCommentCount(article.ID, 1); err != nil {
		logger.Warnf("更新文章评论数失败, articleID: %d, err: %v", article.ID, err)
	}

	go func() {
		defer func() { _ = recover() }()
		s.sendEmailNotifications(comment, article)
	}()

	return comment.ID, nil
}

// sendEmailNotifications 发送邮件通知
func (s *commentService) sendEmailNotifications(comment *entity.Comment, article *entity.Article) {
	adminEmail := s.config.Email.FromEmail
	logger.Infof("开始发送邮件通知，博主邮箱: %s", adminEmail)

	// 如果是回复评论，发送邮件给「被回复的评论」的作者（优先使用 ReplyToID，无则回退到 ParentID 兼容老数据）
	var targetCommentID *uint
	switch {
	case comment.ReplyToID != nil && *comment.ReplyToID > 0:
		targetCommentID = comment.ReplyToID
	case comment.ParentID != nil && *comment.ParentID > 0:
		targetCommentID = comment.ParentID
	}

	if targetCommentID != nil && *targetCommentID > 0 {
		targetComment, err := s.commentRepo.FindByID(*targetCommentID)
		if err != nil || targetComment == nil {
			logger.Warnf("获取被回复评论失败, targetID: %d, err: %v", *targetCommentID, err)
		} else if targetComment.Email != "" && targetComment.Email != comment.Email {
			err = s.emailSvc.SendReplyNotification(
				targetComment.Email,
				comment.Nickname,
				article.Title,
				article.Slug,
				comment.Content,
			)
			if err != nil {
				logger.Warnf("发送回复通知邮件失败, to: %s, err: %v", targetComment.Email, err)
			} else {
				logger.Infof("回复通知邮件发送成功, to: %s", targetComment.Email)
			}
		}
	}

	if comment.ParentID == nil {
		if adminEmail != "" {
			err := s.emailSvc.SendCommentNotification(
				adminEmail,
				comment.Nickname,
				article.Title,
				article.Slug,
				comment.Content,
			)
			if err != nil {
				logger.Warnf("发送评论通知邮件失败, to: %s, err: %v", adminEmail, err)
			} else {
				logger.Infof("评论通知邮件发送成功, to: %s", adminEmail)
			}
		} else {
			logger.Warn("博主邮箱为空，无法发送邮件通知")
		}
	}
}

// GetAdminCommentList 获取评论列表（后台）
func (s *commentService) GetAdminCommentList(req *request.CommentListRequest) (*response.PageResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := s.commentRepo.AdminList((req.Page-1)*req.PageSize, req.PageSize, req.Status)
	if err != nil {
		return nil, fmt.Errorf("获取后台评论列表失败, %w", err)
	}

	respList := make([]response.CommentResponse, 0, len(list))
	for _, c := range list {
		respList = append(respList, s.convertToCommentResponse(c))
	}

	return response.NewPageResponse(respList, total, req.Page, req.PageSize), nil
}

// UpdateCommentStatus 更新评论状态
func (s *commentService) UpdateCommentStatus(id uint, status string) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询评论失败, %w", err)
	}
	if comment == nil {
		return bizerrors.New(bizerrors.CodeCommentNotFound, bizerrors.GetMessage(bizerrors.CodeCommentNotFound))
	}

	oldStatus := comment.Status
	comment.Status = status
	if err := s.commentRepo.Update(comment); err != nil {
		return fmt.Errorf("更新评论状态失败, %w", err)
	}

	// 同步文章评论数：仅统计 approved 评论
	if oldStatus != status {
		var delta int64
		if status == "approved" && oldStatus != "approved" {
			delta = 1
		} else if oldStatus == "approved" && status != "approved" {
			delta = -1
		}
		if delta != 0 {
			if err := s.articleRepo.UpdateCommentCount(comment.ArticleID, delta); err != nil {
				logger.Warnf("更新文章评论数失败, articleID: %d, err: %v", comment.ArticleID, err)
			}
		}
	}
	return nil
}

// DeleteComment 删除评论
func (s *commentService) DeleteComment(id uint) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询评论失败, %w", err)
	}
	if comment == nil {
		return bizerrors.New(bizerrors.CodeCommentNotFound, bizerrors.GetMessage(bizerrors.CodeCommentNotFound))
	}

	if err := s.commentRepo.Delete(id); err != nil {
		return fmt.Errorf("删除评论失败, %w", err)
	}

	// 同步文章评论数（仅 approved 评论计入）
	if comment.Status == "approved" {
		if err := s.articleRepo.UpdateCommentCount(comment.ArticleID, -1); err != nil {
			logger.Warnf("更新文章评论数失败, articleID: %d, err: %v", comment.ArticleID, err)
		}
	}
	return nil
}

// BatchDeleteComments 批量删除评论
func (s *commentService) BatchDeleteComments(ids []uint) error {
	// 先查询待删评论，用于同步文章评论数
	var commentsToDelete []*entity.Comment
	for _, id := range ids {
		comment, err := s.commentRepo.FindByID(id)
		if err != nil {
			return fmt.Errorf("查询评论失败, %w", err)
		}
		if comment != nil {
			commentsToDelete = append(commentsToDelete, comment)
		}
	}

	if err := s.commentRepo.BatchDelete(ids); err != nil {
		return fmt.Errorf("批量删除评论失败, %w", err)
	}

	// 按文章分组统计减少的 approved 评论数
	deltaByArticle := make(map[uint]int64)
	for _, c := range commentsToDelete {
		if c.Status == "approved" {
			deltaByArticle[c.ArticleID]++
		}
	}
	for articleID, delta := range deltaByArticle {
		if err := s.articleRepo.UpdateCommentCount(articleID, -delta); err != nil {
			logger.Warnf("更新文章评论数失败, articleID: %d, err: %v", articleID, err)
		}
	}
	return nil
}

// LikeComment 点赞评论（事务 + 唯一约束防重复）
func (s *commentService) LikeComment(commentID uint, visitorIP string) error {
	// 检查评论是否存在
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return fmt.Errorf("查询评论失败, %w", err)
	}
	if comment == nil {
		return bizerrors.New(bizerrors.CodeCommentNotFound, bizerrors.GetMessage(bizerrors.CodeCommentNotFound))
	}

	// 事务：插入点赞日志（防重复）+ 增加点赞数
	return s.commentRepo.LikeWithLog(commentID, visitorIP)
}

// UnlikeComment 取消点赞评论（事务）
func (s *commentService) UnlikeComment(commentID uint, visitorIP string) error {
	// 检查评论是否存在
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return fmt.Errorf("查询评论失败, %w", err)
	}
	if comment == nil {
		return bizerrors.New(bizerrors.CodeCommentNotFound, bizerrors.GetMessage(bizerrors.CodeCommentNotFound))
	}

	// 事务：删除点赞日志 + 减少点赞数
	return s.commentRepo.UnlikeWithLog(commentID, visitorIP)
}
