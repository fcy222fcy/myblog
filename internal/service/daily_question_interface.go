package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
)

// DailyQuestionService 每日一问服务接口
type DailyQuestionService interface {
	// GetLatestQuestion 获取最新问题
	GetLatestQuestion() (*response.DailyQuestionResponse, error)

	// GetAllPublishedQuestions 获取所有已发布问题列表
	GetAllPublishedQuestions() ([]*response.DailyQuestionBriefResponse, error)

	// GetQuestionByDate 获取指定日期问题
	GetQuestionByDate(date string) (*response.DailyQuestionResponse, error)

	// GetPreviousQuestion 获取前一天的问题
	GetPreviousQuestion(date string) (*response.DailyQuestionResponse, error)

	// GetNextQuestion 获取后一天的问题
	GetNextQuestion(date string) (*response.DailyQuestionResponse, error)

	// LikeQuestion 问题点赞（按 IP 防重复）
	LikeQuestion(id uint, visitorIP string) (int64, error)

	// GetAdminQuestionList 获取问题列表（后台）
	GetAdminQuestionList(req *request.DailyQuestionListRequest) (*response.PageResponse, error)

	// CreateQuestion 创建问题
	CreateQuestion(req *request.CreateDailyQuestionRequest) (uint, error)

	// UpdateQuestion 更新问题
	UpdateQuestion(id uint, req *request.UpdateDailyQuestionRequest) error

	// DeleteQuestion 删除问题
	DeleteQuestion(id uint) error

	// UpdateQuestionStatus 更新问题状态
	UpdateQuestionStatus(id uint, status int) error
}
