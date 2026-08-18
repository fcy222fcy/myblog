package repository

import (
	"blog/internal/model/entity"
	"strings"
	"time"

	"gorm.io/gorm"
)

// articleRepository 文章数据访问实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章数据访问
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// FindByID 根据 ID 查找文章
func (r *articleRepository) FindByID(id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.Preload("Category").Preload("Tags").First(&article, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// FindBySlug 根据 slug 查找已发布文章
func (r *articleRepository) FindBySlug(slug string) (*entity.Article, error) {
	var article entity.Article
	err := r.db.Preload("Category").Preload("Tags").
		Where("slug = ? AND status = ?", slug, entity.ArticleStatusPublished).
		First(&article).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// Create 创建文章
func (r *articleRepository) Create(article *entity.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *articleRepository) Update(article *entity.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章（硬删除，级联清理评论、点赞记录、标签关联）
func (r *articleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 查询该文章的所有评论 ID
		var commentIDs []uint
		if err := tx.Model(&entity.Comment{}).
			Where("article_id = ?", id).
			Pluck("id", &commentIDs).Error; err != nil {
			return err
		}

		// 删除评论点赞记录
		if len(commentIDs) > 0 {
			if err := tx.Where("comment_id IN ?", commentIDs).
				Delete(&entity.CommentLikeLog{}).Error; err != nil {
				return err
			}
		}

		// 删除评论
		if err := tx.Where("article_id = ?", id).
			Delete(&entity.Comment{}).Error; err != nil {
			return err
		}

		// 清理文章-标签关联
		if err := tx.Where("article_id = ?", id).
			Delete(&entity.ArticleTag{}).Error; err != nil {
			return err
		}

		// 删除文章
		return tx.Delete(&entity.Article{}, id).Error
	})
}

// ListPublished 已发布文章列表（前台）
func (r *articleRepository) ListPublished(offset, limit int, categoryId, tagId uint, keyword string) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{}).Where("status = ?", entity.ArticleStatusPublished)
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if categoryId > 0 {
		query = query.Where("category_id = ?", categoryId)
	}
	if tagId > 0 {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", tagId)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Category").Preload("Tags").
		Offset(offset).Limit(limit).
		Order("is_top DESC, created_at DESC").
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ListAll 所有文章列表（后台）
func (r *articleRepository) ListAll(offset, limit int, status, keyword string, categoryID uint) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Category").Preload("Tags").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// IncrementViewCount 增加浏览量
func (r *articleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&entity.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateCommentCount 调整文章评论数（delta 为增量，可为负数，带下限保护）
func (r *articleRepository) UpdateCommentCount(articleID uint, delta int64) error {
	if delta == 0 {
		return nil
	}
	if delta > 0 {
		return r.db.Model(&entity.Article{}).Where("id = ?", articleID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", delta)).Error
	}
	abs := -delta
	return r.db.Model(&entity.Article{}).Where("id = ?", articleID).
		UpdateColumn("comment_count",
			gorm.Expr("CASE WHEN comment_count >= ? THEN comment_count - ? ELSE comment_count END", abs, abs)).Error
}

// BatchDelete 批量删除（逐个处理级联清理）
func (r *articleRepository) BatchDelete(ids []uint) error {
	for _, id := range ids {
		if err := r.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

// Count 统计文章数量
func (r *articleRepository) Count(status string) (int64, error) {
	var total int64
	query := r.db.Model(&entity.Article{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	return total, err
}

// SumViewCount 统计总浏览量
func (r *articleRepository) SumViewCount() (int64, error) {
	var sum int64
	err := r.db.Model(&entity.Article{}).Select("COALESCE(SUM(view_count), 0)").Scan(&sum).Error
	return sum, err
}

// FindByCategoryID 根据分类ID查找文章
func (r *articleRepository) FindByCategoryID(categoryID uint, offset, limit int) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{}).Where("category_id = ? AND status = ?", categoryID, entity.ArticleStatusPublished)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Category").Preload("Tags").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&articles).Error
	return articles, total, err
}

// FindByTagID 根据标签ID查找文章
func (r *articleRepository) FindByTagID(tagID uint, offset, limit int) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{}).
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Where("article_tags.tag_id = ? AND articles.status = ?", tagID, entity.ArticleStatusPublished)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Category").Preload("Tags").
		Offset(offset).Limit(limit).
		Order("articles.created_at DESC").
		Find(&articles).Error
	return articles, total, err
}

// GetArchives 获取文章归档
func (r *articleRepository) GetArchives() ([]map[string][]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.Where("status = ?", entity.ArticleStatusPublished).
		Order("created_at DESC").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	// 按年份分组
	yearMap := make(map[string][]*entity.Article)
	for _, article := range articles {
		year := article.CreatedAt.Format("2006")
		yearMap[year] = append(yearMap[year], article)
	}

	var result []map[string][]*entity.Article
	for year, arts := range yearMap {
		result = append(result, map[string][]*entity.Article{year: arts})
	}
	return result, nil
}

// GetRecent 获取最近文章
func (r *articleRepository) GetRecent(limit int) ([]entity.Article, error) {
	var articles []entity.Article
	err := r.db.
		Where("status = ?", entity.ArticleStatusPublished).
		Order("created_at DESC").
		Limit(limit).
		Preload("Category").
		Preload("Tags").
		Find(&articles).Error
	return articles, err
}

// escapeLikePattern 转义 SQL LIKE 通配符（% 和 _），防止用户输入被当作通配符
func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// Search 搜索文章（标题、内容、摘要、标签名、分类名模糊搜索），按相关性排序
func (r *articleRepository) Search(keyword string, offset, limit int) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	// 按空格拆分关键词，每个词独立匹配（OR 逻辑）
	words := strings.Fields(keyword)
	if len(words) == 0 {
		return nil, 0, nil
	}

	// 构建每个词的 LIKE 条件
	// whereArgs: 用于 WHERE 和子查询的参数（title, content, summary, category, tag）
	// scoreArgs: 用于相关性评分 CASE WHEN 的参数（title, summary, tag）
	var titleConds, contentConds, summaryConds, catConds, tagConds []string
	var titleScoreConds, summaryScoreConds, tagScoreConds []string
	var whereArgs, scoreArgs []interface{}

	for _, w := range words {
		likeVal := "%" + escapeLikePattern(w) + "%"
		titleConds = append(titleConds, "a.title LIKE ?")
		contentConds = append(contentConds, "a.content LIKE ?")
		summaryConds = append(summaryConds, "a.summary LIKE ?")
		catConds = append(catConds, "c.name LIKE ?")
		tagConds = append(tagConds, "t.name LIKE ?")
		whereArgs = append(whereArgs, likeVal, likeVal, likeVal, likeVal, likeVal)

		titleScoreConds = append(titleScoreConds, "a.title LIKE ?")
		summaryScoreConds = append(summaryScoreConds, "a.summary LIKE ?")
		tagScoreConds = append(tagScoreConds, "t.name LIKE ?")
		scoreArgs = append(scoreArgs, likeVal, likeVal, likeVal)
	}

	titleExpr := "(" + strings.Join(titleConds, " OR ") + ")"
	contentExpr := "(" + strings.Join(contentConds, " OR ") + ")"
	summaryExpr := "(" + strings.Join(summaryConds, " OR ") + ")"
	catExpr := "(" + strings.Join(catConds, " OR ") + ")"
	tagExpr := "(" + strings.Join(tagConds, " OR ") + ")"
	whereCondition := titleExpr + " OR " + contentExpr + " OR " + summaryExpr + " OR " + catExpr + " OR " + tagExpr

	// 相关性评分 CASE WHEN：标题(10) > 摘要(5) > 标签(3) > 内容(1)
	titleScoreExpr := "(" + strings.Join(titleScoreConds, " OR ") + ")"
	summaryScoreExpr := "(" + strings.Join(summaryScoreConds, " OR ") + ")"
	tagScoreExpr := "(" + strings.Join(tagScoreConds, " OR ") + ")"
	relevanceScore := "CASE WHEN " + titleScoreExpr + " THEN 10 WHEN " + summaryScoreExpr + " THEN 5 WHEN " + tagScoreExpr + " THEN 3 ELSE 1 END"

	// 子查询：匹配标签名的文章 ID（只用 tag 的 whereArgs）
	tagStart := len(words) * 4 // tag 条件在 whereArgs 中的起始位置
	tagSubQuery := r.db.Model(&entity.Article{}).
		Select("DISTINCT articles.id").
		Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Joins("JOIN tags t ON t.id = article_tags.tag_id").
		Where("articles.status = ?", "published").
		Where(tagExpr, whereArgs[tagStart:]...)

	// 主查询 + COUNT 共用的 WHERE 条件
	mainWhere := append(whereArgs, tagSubQuery)

	query := r.db.Model(&entity.Article{}).
		Table("articles a").
		Joins("LEFT JOIN categories c ON c.id = a.category_id").
		Joins("LEFT JOIN article_tags at2 ON at2.article_id = a.id").
		Joins("LEFT JOIN tags t ON t.id = at2.tag_id").
		Where("a.status = ?", "published").
		Where("("+whereCondition+") OR a.id IN (?)", mainWhere...).
		Group("a.id").
		Select("a.*, " + relevanceScore + " AS relevance_score")

	// COUNT 查询（不带排序和分页）
	countQuery := r.db.Model(&entity.Article{}).
		Table("articles a").
		Joins("LEFT JOIN categories c ON c.id = a.category_id").
		Joins("LEFT JOIN article_tags at2 ON at2.article_id = a.id").
		Joins("LEFT JOIN tags t ON t.id = at2.tag_id").
		Where("a.status = ?", "published").
		Where("("+whereCondition+") OR a.id IN (?)", mainWhere...).
		Group("a.id")

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Category").Preload("Tags").
		Order("relevance_score DESC, a.is_top DESC, a.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetDB 获取数据库实例
func (r *articleRepository) GetDB() *gorm.DB {
	return r.db
}

// UpdateTags 更新文章标签关联
func (r *articleRepository) UpdateTags(article *entity.Article, tags []entity.Tag) error {
	return r.db.Model(article).Association("Tags").Replace(tags)
}

func (r *articleRepository) ListScheduledAfter(now time.Time) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.
		Where("status = ? AND scheduled_at > ?", entity.ArticleStatusScheduled, now).
		Find(&articles).Error
	return articles, err
}

func (r *articleRepository) PublishScheduledArticle(id uint, now time.Time) (bool, error) {
	result := r.db.Model(&entity.Article{}).
		Where("id = ? AND status = ? AND scheduled_at <= ?", id, entity.ArticleStatusScheduled, now).
		Updates(map[string]interface{}{
			"status":       entity.ArticleStatusPublished,
			"scheduled_at": nil,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *articleRepository) PublishDueScheduledArticles(now time.Time) (int64, error) {
	result := r.db.Model(&entity.Article{}).
		Where("status = ? AND scheduled_at <= ?", entity.ArticleStatusScheduled, now).
		Updates(map[string]interface{}{
			"status":       entity.ArticleStatusPublished,
			"scheduled_at": nil,
		})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
