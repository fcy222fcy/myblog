package response

import "time"

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID           uint             `json:"id"`
	Title        string           `json:"title"`
	Slug         string           `json:"slug"`
	Content      string           `json:"content,omitempty"`
	Summary      string           `json:"summary"`
	Cover        string           `json:"cover"`
	Category     CategoryResponse `json:"category"`
	Tags         []TagResponse    `json:"tags"`
	ViewCount    int64            `json:"view_count"`
	CommentCount int64            `json:"comment_count"`
	Status       string           `json:"status"`
	IsTop        bool             `json:"is_top"`
	ReadingTime  int              `json:"reading_time"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// ArticleDetailResponse 文章详情响应
type ArticleDetailResponse struct {
	ID           uint             `json:"id"`
	Title        string           `json:"title"`
	Slug         string           `json:"slug"`
	Content      string           `json:"content"`
	Summary      string           `json:"summary"`
	Cover        string           `json:"cover"`
	Category     CategoryResponse `json:"category"`
	Tags         []TagResponse    `json:"tags"`
	ViewCount    int64            `json:"view_count"`
	CommentCount int64            `json:"comment_count"`
	Status       string           `json:"status"`
	IsTop        bool             `json:"is_top"`
	ReadingTime  int              `json:"reading_time"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// AdminArticleDetailResponse 后台文章详情响应
type AdminArticleDetailResponse struct {
	ID             uint       `json:"id"`
	Title          string     `json:"title"`
	Slug           string     `json:"slug"`
	Content        string     `json:"content"`
	Summary        string     `json:"summary"`
	Cover          string     `json:"cover"`
	Status         string     `json:"status"`
	IsTop          bool       `json:"is_top"`
	CategoryID     uint       `json:"category_id"`
	TagIDs         []uint     `json:"tag_ids"`
	ScheduledAt    *time.Time `json:"scheduled_at"`
	SEOTitle       string     `json:"seo_title"`
	SEODescription string     `json:"seo_description"`
	SEOKeywords    string     `json:"seo_keywords"`
}

// ArticleListResponse 文章列表响应（简化版）
type ArticleListResponse struct {
	ID           uint             `json:"id"`
	Title        string           `json:"title"`
	Slug         string           `json:"slug"`
	Summary      string           `json:"summary"`
	Cover        string           `json:"cover"`
	Category     CategoryResponse `json:"category"`
	Tags         []TagResponse    `json:"tags"`
	ViewCount    int64            `json:"view_count"`
	CommentCount int64            `json:"comment_count"`
	ReadingTime  int              `json:"reading_time"`
	CreatedAt    time.Time        `json:"created_at"`
}

// ArchiveResponse 归档响应
type ArchiveResponse struct {
	Year     int                      `json:"year"`
	Articles []ArchiveArticleResponse `json:"articles"`
}

// ArchiveArticleResponse 归档文章响应
type ArchiveArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// SearchArticleResponse 搜索结果文章响应（带关键词上下文片段）
type SearchArticleResponse struct {
	ID            uint             `json:"id"`
	Title         string           `json:"title"`
	Slug          string           `json:"slug"`
	Summary       string           `json:"summary"`
	SearchSnippet string           `json:"search_snippet,omitempty"`
	Cover         string           `json:"cover"`
	Category      CategoryResponse `json:"category"`
	Tags          []TagResponse    `json:"tags"`
	ViewCount     int64            `json:"view_count"`
	CommentCount  int64            `json:"comment_count"`
	ReadingTime   int              `json:"reading_time"`
	CreatedAt     time.Time        `json:"created_at"`
}
