package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
)

// ArticleService 文章服务接口
type ArticlePublishScheduler interface {
	RefreshArticle(article *entity.Article) error
	UnscheduleArticle(id uint)
}

type ArticleService interface {
	SetPublishScheduler(scheduler ArticlePublishScheduler)
	HandleScheduledArticlePublished(article *entity.Article)

	// GetArticleList 获取文章列表（前台）
	GetArticleList(req *request.ArticleListRequest) (*response.PageResponse, error)

	// GetArticleDetail 获取文章详情（前台）
	GetArticleDetail(slug string) (*response.ArticleDetailResponse, error)

	// GetArticleArchives 获取文章归档
	GetArticleArchives() ([]response.ArchiveResponse, error)

	// GetAdminArticleList 获取文章列表（后台）
	GetAdminArticleList(req *request.ArticleListRequest) (*response.PageResponse, error)

	// GetAdminArticleDetail 获取文章详情（后台）
	GetAdminArticleDetail(id uint) (*response.AdminArticleDetailResponse, error)

	// CreateArticle 创建文章
	CreateArticle(req *request.CreateArticleRequest) (uint, error)

	// UpdateArticle 更新文章
	UpdateArticle(id uint, req *request.UpdateArticleRequest) error

	// UpdateArticleStatus 更新文章状态（列表页快捷切换）
	UpdateArticleStatus(id uint, status string) error

	// DeleteArticle 删除文章
	DeleteArticle(id uint) error

	// BatchDeleteArticles 批量删除文章
	BatchDeleteArticles(ids []uint) error

	// Search 搜索文章（标题、内容、摘要模糊搜索）
	Search(keyword string, page, pageSize int) (*response.PageResponse, error)
}
