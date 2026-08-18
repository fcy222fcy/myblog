package repository

import (
	"blog/internal/model/entity"
	"time"

	"gorm.io/gorm"
)

// ArticleRepository 文章数据访问接口
type ArticleRepository interface {
	// FindByID 根据 ID 查找文章
	FindByID(id uint) (*entity.Article, error)

	// FindBySlug 根据 slug 查找文章
	FindBySlug(slug string) (*entity.Article, error)

	// Create 创建文章
	Create(article *entity.Article) error

	// Update 更新文章
	Update(article *entity.Article) error

	// Delete 删除文章（硬删除，级联清理评论和标签关联）
	Delete(id uint) error

	// ListPublished 已发布文章列表（前台）
	ListPublished(offset, limit int, categoryId, tagId uint, keyword string) ([]*entity.Article, int64, error)

	// ListAll 所有文章列表（后台）
	ListAll(offset, limit int, status, keyword string, categoryID uint) ([]*entity.Article, int64, error)

	// IncrementViewCount 增加浏览量
	IncrementViewCount(id uint) error

	// UpdateCommentCount 调整文章评论数（delta 为增量，可为负数）
	UpdateCommentCount(articleID uint, delta int64) error

	// BatchDelete 批量删除
	BatchDelete(ids []uint) error

	// Count 统计文章数量
	Count(status string) (int64, error)

	// SumViewCount 统计总浏览量
	SumViewCount() (int64, error)

	// FindByCategoryID 根据分类ID查找文章
	FindByCategoryID(categoryID uint, offset, limit int) ([]*entity.Article, int64, error)

	// FindByTagID 根据标签ID查找文章
	FindByTagID(tagID uint, offset, limit int) ([]*entity.Article, int64, error)

	// GetArchives 获取文章归档（按年份分组）
	GetArchives() ([]map[string][]*entity.Article, error)

	// GetRecent 获取最近文章
	GetRecent(limit int) ([]entity.Article, error)

	// Search 搜索文章（标题、内容、摘要模糊搜索）
	Search(keyword string, offset, limit int) ([]*entity.Article, int64, error)

	// UpdateTags 更新文章标签关联
	UpdateTags(article *entity.Article, tags []entity.Tag) error

	// GetDB 获取数据库实例（用于事务）
	ListScheduledAfter(now time.Time) ([]*entity.Article, error)

	PublishScheduledArticle(id uint, now time.Time) (bool, error)

	PublishDueScheduledArticles(now time.Time) (int64, error)

	GetDB() *gorm.DB
}
