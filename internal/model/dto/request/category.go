package request

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Slug        string `json:"slug" binding:"max=50"`
	Description string `json:"description" binding:"max=200"`
	Icon        string `json:"icon" binding:"max=10"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string  `json:"name" binding:"min=1,max=50"`
	Slug        *string `json:"slug" binding:"max=50"`
	Description *string `json:"description" binding:"max=200"`
	Icon        *string `json:"icon" binding:"max=10"`
	SortOrder   int     `json:"sort_order"`
}
