package entity

// Comment 评论
type Comment struct {
	BaseEntity
	Content string `gorm:"type:text;not null" json:"content"`

	// 注册用户评论
	UserID *uint `gorm:"index" json:"user_id"`
	User   *User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 兼容字段（内联存储访客信息）
	Nickname string `gorm:"type:varchar(50)" json:"nickname"`
	Email    string `gorm:"type:varchar(100)" json:"email"`
	Website  string `gorm:"type:varchar(200)" json:"website"`
	Avatar   string `gorm:"type:varchar(500)" json:"avatar"`

	ArticleID       uint      `gorm:"index" json:"article_id"`
	Article         Article   `gorm:"foreignKey:ArticleID" json:"-"`
	ParentID        *uint     `gorm:"index" json:"parent_id"`
	ReplyToID       *uint     `gorm:"index" json:"reply_to_id"`
	ReplyTo         *Comment  `gorm:"foreignKey:ReplyToID" json:"-"`
	ReplyToNickname string    `gorm:"-" json:"reply_to_nickname"`
	Replies         []Comment `gorm:"-" json:"replies,omitempty"`
	LikeCount       int       `gorm:"default:0" json:"like_count"`
	Status          string    `gorm:"type:varchar(20);default:pending;index" json:"status"` // pending: 待审核 approved: 已通过 rejected: 已拒绝
	IP              string    `gorm:"type:varchar(50)" json:"ip"`
	UserAgent       string    `gorm:"type:varchar(500)" json:"user_agent"`

	// 解析后的客户端信息（优先前端检测，后端 UA 解析兜底）
	OS             string `gorm:"type:varchar(50)" json:"os"`
	OSVersion      string `gorm:"type:varchar(50)" json:"os_version"`
	Browser        string `gorm:"type:varchar(50)" json:"browser"`
	BrowserVersion string `gorm:"type:varchar(50)" json:"browser_version"`
}

// TableName 表名
func (Comment) TableName() string {
	return "comments"
}
