package request

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	PageRequest
	Status    string `json:"status" form:"status"`
	Keyword   string `json:"keyword" form:"keyword"`
	ArticleID uint   `json:"article_id" form:"article_id"`
	SortBy    string `json:"sort_by" form:"sort_by"` // asc / desc / hot
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	ArticleID   uint   `json:"article_id"`
	ArticleSlug string `json:"article_slug"`
	Content     string `json:"content" binding:"required,min=1,max=1000"`
	Nickname    string `json:"nickname" binding:"required,min=1,max=50"`
	Email       string `json:"email" binding:"omitempty,email,max=100"`
	Website     string `json:"website" binding:"max=200"`
	ParentID    *uint  `json:"parent_id"`
	ReplyToID   *uint  `json:"reply_to_id"`

	// 客户端信息（优先前端 JS 检测，解决 Win11 UA 无法区分等问题；不传则后端解析 UA 兜底）
	OS             string `json:"os"`
	OSVersion      string `json:"os_version"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
}

// UpdateCommentStatusRequest 更新评论状态请求
type UpdateCommentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}
