package response

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	SortOrder    int    `json:"sort_order"`
	ArticleCount int64  `json:"article_count"`
}
