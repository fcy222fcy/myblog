package entity

// AboutPage 关于页面
type AboutPage struct {
	BaseEntity
	Title       string `gorm:"type:varchar(200)" json:"title"`    // 页面标题
	Subtitle    string `gorm:"type:varchar(500)" json:"subtitle"` // 副标题
	Bio         string `gorm:"type:text" json:"bio"`              // 简介
	Skills      string `gorm:"type:text" json:"skills"`           // 技能标签（JSON 数组）
	AboutMe     string `gorm:"type:text" json:"about_me"`         // 关于我详细信息（JSON）
	AboutSite   string `gorm:"type:text" json:"about_site"`       // 关于网站信息（JSON）
	Projects    string `gorm:"type:text" json:"projects"`         // 项目列表（JSON）
	ContactInfo string `gorm:"type:text" json:"contact_info"`     // 联系信息（JSON）
	SiteHistory string `gorm:"type:text" json:"site_history"`     // 建站历程（JSON 数组）
}

// TableName 表名
func (AboutPage) TableName() string {
	return "about_pages"
}

// SkillTag 技能标签
type SkillTag struct {
	Name string `json:"name"`
}

// AboutMeItem 关于我信息项
type AboutMeItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
}

// ProjectItem 项目项
type ProjectItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
}

// ContactItem 联系方式项
type ContactItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
	Url   string `json:"url"`
}

// SiteHistoryItem 建站历程项
type SiteHistoryItem struct {
	Date    string `json:"date"`
	Content string `json:"content"`
	Icon    string `json:"icon"`
}
