package entity

// AuditLog 操作审计日志
type AuditLog struct {
	BaseEntity
	OperatorID   uint   `gorm:"index" json:"operator_id"`                      // 操作人ID
	OperatorName string `gorm:"type:varchar(50)" json:"operator_name"`         // 操作人名称
	Action       string `gorm:"type:varchar(50);index;not null" json:"action"` // 操作类型：create/update/delete/approve/reject
	TargetType   string `gorm:"type:varchar(50);index" json:"target_type"`     // 目标类型：article/comment/category/tag/link
	TargetID     uint   `gorm:"index" json:"target_id"`                        // 目标ID
	TargetTitle  string `gorm:"type:varchar(200)" json:"target_title"`         // 目标标题（冗余字段，方便查看）
	Detail       string `gorm:"type:text" json:"detail"`                       // 操作详情（JSON 格式）
	IP           string `gorm:"type:varchar(50)" json:"ip"`                    // 操作IP
	UserAgent    string `gorm:"type:varchar(500)" json:"user_agent"`           // User-Agent
}

// TableName 表名
func (AuditLog) TableName() string {
	return "audit_logs"
}
