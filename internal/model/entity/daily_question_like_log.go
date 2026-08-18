package entity

import "time"

// DailyQuestionLikeLog 每日一问点赞记录
type DailyQuestionLikeLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	QuestionID uint      `gorm:"uniqueIndex:idx_question_ip;not null" json:"question_id"`
	VisitorIP  string    `gorm:"type:varchar(50);uniqueIndex:idx_question_ip;not null" json:"visitor_ip"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 表名
func (DailyQuestionLikeLog) TableName() string {
	return "daily_question_like_logs"
}
