package request

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// DailyQuestionListRequest 每日一问列表请求
type DailyQuestionListRequest struct {
	PageRequest
	Status  *int   `json:"status" form:"status"`
	Keyword string `json:"keyword" form:"keyword"`
	Date    string `json:"date" form:"date"`
}

// CreateDailyQuestionRequest 创建每日一问请求
type CreateDailyQuestionRequest struct {
	Question string `json:"question" binding:"required,min=1,max=500"`
	Answer   string `json:"answer" binding:"required,max=2000"`
	Date     string `json:"date" binding:"required,datetime=2006-01-02"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// UpdateDailyQuestionRequest 更新每日一问请求
type UpdateDailyQuestionRequest struct {
	Question string  `json:"question" binding:"min=1"`
	Answer   *string `json:"answer"`
	Date     string  `json:"date" binding:"omitempty,datetime=2006-01-02"`
	Status   *int    `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// ValidateDate 验证日期格式是否为 YYYY-MM-DD
func ValidateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return true // 空值由 required 验证处理
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
