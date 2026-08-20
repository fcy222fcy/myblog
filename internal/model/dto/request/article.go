package request

import "time"

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	PageRequest
	Category uint   `json:"category_id" form:"category_id"`
	Tag      uint   `json:"tag" form:"tag"`
	Keyword  string `json:"keyword" form:"keyword"`
	Status   string `json:"status" form:"status"`
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title          string     `json:"title" binding:"required,min=1,max=200"`
	Content        string     `json:"content" binding:"required,max=50000"`
	Summary        string     `json:"summary" binding:"max=500"`
	Cover          string     `json:"cover" binding:"max=500"`
	CategoryID     uint       `json:"category_id" binding:"required"`
	TagIDs         []uint     `json:"tag_ids"`
	Status         string     `json:"status" binding:"omitempty,oneof=published draft scheduled"`
	IsTop          bool       `json:"is_top"`
	Slug           string     `json:"slug" binding:"max=200"`
	ScheduledAt    *time.Time `json:"scheduled_at"`
	SEOTitle       string     `json:"seo_title" binding:"max=200"`
	SEODescription string     `json:"seo_description" binding:"max=500"`
	SEOKeywords    string     `json:"seo_keywords" binding:"max=300"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title          string     `json:"title" binding:"min=1,max=200"`
	Content        string     `json:"content"`
	Summary        *string    `json:"summary"`
	Cover          *string    `json:"cover"`
	CategoryID     uint       `json:"category_id"`
	TagIDs         []uint     `json:"tag_ids"`
	Status         string     `json:"status" binding:"omitempty,oneof=published draft scheduled"`
	IsTop          bool       `json:"is_top"`
	Slug           string     `json:"slug" binding:"max=200"`
	ScheduledAt    *time.Time `json:"scheduled_at"`
	SEOTitle       *string    `json:"seo_title" binding:"max=200"`
	SEODescription *string    `json:"seo_description" binding:"max=500"`
	SEOKeywords    *string    `json:"seo_keywords" binding:"max=300"`
}
