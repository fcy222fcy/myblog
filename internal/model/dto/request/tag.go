package request

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
	Slug string `json:"slug" binding:"max=50"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name string  `json:"name" binding:"min=1,max=50"`
	Slug *string `json:"slug" binding:"max=50"`
}
