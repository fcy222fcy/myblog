package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/logger"
	"fmt"
	"strings"
	"time"
)

// dailyQuestionService 每日一问服务实现
type dailyQuestionService struct {
	dailyQuestionRepo repository.DailyQuestionRepository
}

// NewDailyQuestionService 创建每日一问服务
func NewDailyQuestionService(dailyQuestionRepo repository.DailyQuestionRepository) DailyQuestionService {
	return &dailyQuestionService{dailyQuestionRepo: dailyQuestionRepo}
}

// GetAllPublishedQuestions 获取所有已发布问题列表
func (s *dailyQuestionService) GetAllPublishedQuestions() ([]*response.DailyQuestionBriefResponse, error) {
	list, err := s.dailyQuestionRepo.GetAllPublished()
	if err != nil {
		return nil, fmt.Errorf("获取已发布问题列表失败, %w", err)
	}

	var result []*response.DailyQuestionBriefResponse
	for _, q := range list {
		result = append(result, &response.DailyQuestionBriefResponse{
			ID:       q.ID,
			Question: q.Question,
			Date:     q.Date,
		})
	}
	return result, nil
}

// GetLatestQuestion 获取最新问题
func (s *dailyQuestionService) GetLatestQuestion() (*response.DailyQuestionResponse, error) {
	question, err := s.dailyQuestionRepo.GetLatest()
	if err != nil {
		return nil, fmt.Errorf("获取最新问题失败, %w", err)
	}
	if question == nil {
		return nil, bizerrors.New(bizerrors.CodeDailyQuestionNotFound, bizerrors.GetMessage(bizerrors.CodeDailyQuestionNotFound))
	}

	_ = s.dailyQuestionRepo.IncrementViewCount(question.ID)
	return s.toResponse(question), nil
}

// GetQuestionByDate 根据日期获取问题（仅允许查看今天及以前的已发布题目，未来日期不提前公开）
func (s *dailyQuestionService) GetQuestionByDate(date string) (*response.DailyQuestionResponse, error) {
	question, err := s.dailyQuestionRepo.FindByDate(date)
	if err != nil {
		return nil, fmt.Errorf("根据日期查询问题失败, %w", err)
	}
	if question == nil {
		return nil, bizerrors.New(bizerrors.CodeDailyQuestionNotFound, bizerrors.GetMessage(bizerrors.CodeDailyQuestionNotFound))
	}
	// 仅限制已发布状态
	if question.Status != entity.DailyQuestionStatusPublished {
		return nil, bizerrors.New(bizerrors.CodeDailyQuestionNotFound, bizerrors.GetMessage(bizerrors.CodeDailyQuestionNotFound))
	}
	// 未来日期不提前展示
	if question.Date > time.Now().Format("2006-01-02") {
		return nil, bizerrors.New(bizerrors.CodeDailyQuestionNotFound, bizerrors.GetMessage(bizerrors.CodeDailyQuestionNotFound))
	}

	_ = s.dailyQuestionRepo.IncrementViewCount(question.ID)
	return s.toResponse(question), nil
}

// GetPreviousQuestion 获取前一天的问题
func (s *dailyQuestionService) GetPreviousQuestion(date string) (*response.DailyQuestionResponse, error) {
	question, err := s.dailyQuestionRepo.GetPrevious(date)
	if err != nil {
		return nil, fmt.Errorf("获取前一天问题失败, %w", err)
	}
	if question == nil {
		return nil, bizerrors.New(bizerrors.CodeDailyQuestionNotFound, "没有前一天的问题")
	}

	return s.toResponse(question), nil
}

// GetNextQuestion 获取后一天的问题
func (s *dailyQuestionService) GetNextQuestion(date string) (*response.DailyQuestionResponse, error) {
	question, err := s.dailyQuestionRepo.GetNext(date)
	if err != nil {
		return nil, fmt.Errorf("获取后一天问题失败, %w", err)
	}
	if question == nil {
		return nil, bizerrors.New(bizerrors.CodeDailyQuestionNotFound, "没有后一天的问题")
	}

	return s.toResponse(question), nil
}

// LikeQuestion 问题点赞（按 IP 防重复）
func (s *dailyQuestionService) LikeQuestion(id uint, visitorIP string) (int64, error) {
	count, err := s.dailyQuestionRepo.LikeQuestionWithLog(id, visitorIP)
	if err != nil {
		return 0, err
	}
	logger.Infof("问题点赞成功, id: %d, likeCount: %d", id, count)
	return count, nil
}

// GetAdminQuestionList 获取问题列表（后台）
func (s *dailyQuestionService) GetAdminQuestionList(req *request.DailyQuestionListRequest) (*response.PageResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	status := -1
	if req.Status != nil {
		status = *req.Status
	}
	if status < -1 || status > entity.DailyQuestionStatusScheduled {
		return nil, bizerrors.New(bizerrors.CodeInvalidParams, "无效的问题状态")
	}

	list, total, err := s.dailyQuestionRepo.List((req.Page-1)*req.PageSize, req.PageSize, status, req.Keyword, req.Date)
	if err != nil {
		return nil, fmt.Errorf("获取后台问题列表失败, %w", err)
	}

	var result []response.DailyQuestionResponse
	for _, q := range list {
		result = append(result, *s.toResponse(q))
	}

	return response.NewPageResponse(result, total, req.Page, req.PageSize), nil
}

// CreateQuestion 创建问题
func (s *dailyQuestionService) CreateQuestion(req *request.CreateDailyQuestionRequest) (uint, error) {
	// 验证问题内容
	if strings.TrimSpace(req.Question) == "" {
		return 0, bizerrors.New(bizerrors.CodeInvalidParams, "问题内容不能为空")
	}

	// 验证答案
	if strings.TrimSpace(req.Answer) == "" {
		return 0, bizerrors.New(bizerrors.CodeInvalidParams, "答案不能为空")
	}

	// 验证日期格式
	if req.Date == "" {
		return 0, bizerrors.New(bizerrors.CodeInvalidParams, "日期不能为空")
	}

	existing, _ := s.dailyQuestionRepo.FindByDate(req.Date)
	if existing != nil {
		return 0, bizerrors.New(bizerrors.CodeDailyQuestionDateExists, bizerrors.GetMessage(bizerrors.CodeDailyQuestionDateExists))
	}

	status := entity.DailyQuestionStatusPublished
	if req.Status != nil {
		status = *req.Status
	}
	if status < entity.DailyQuestionStatusDisabled || status > entity.DailyQuestionStatusScheduled {
		return 0, bizerrors.New(bizerrors.CodeInvalidParams, "无效的问题状态")
	}
	question := &entity.DailyQuestion{
		Question: req.Question,
		Answer:   req.Answer,
		Date:     req.Date,
		Status:   status,
	}

	err := s.dailyQuestionRepo.Create(question)
	if err != nil {
		return 0, fmt.Errorf("创建问题失败, %w", err)
	}

	logger.Infof("创建问题成功, id: %d, date: %s", question.ID, question.Date)
	return question.ID, nil
}

// UpdateQuestion 更新问题
func (s *dailyQuestionService) UpdateQuestion(id uint, req *request.UpdateDailyQuestionRequest) error {
	question, err := s.dailyQuestionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询问题失败, %w", err)
	}
	if question == nil {
		return bizerrors.New(bizerrors.CodeDailyQuestionNotFound, bizerrors.GetMessage(bizerrors.CodeDailyQuestionNotFound))
	}

	if req.Question != "" {
		question.Question = req.Question
	}
	if req.Answer != nil {
		question.Answer = *req.Answer
	}
	if req.Date != "" {
		existing, _ := s.dailyQuestionRepo.FindByDate(req.Date)
		if existing != nil && existing.ID != id {
			return bizerrors.New(bizerrors.CodeDailyQuestionDateExists, bizerrors.GetMessage(bizerrors.CodeDailyQuestionDateExists))
		}
		question.Date = req.Date
	}
	if req.Status != nil {
		if *req.Status < entity.DailyQuestionStatusDisabled || *req.Status > entity.DailyQuestionStatusScheduled {
			return bizerrors.New(bizerrors.CodeInvalidParams, "无效的问题状态")
		}
		question.Status = *req.Status
	}

	if err := s.dailyQuestionRepo.Update(question); err != nil {
		return fmt.Errorf("更新问题失败, %w", err)
	}

	logger.Infof("更新问题成功, id: %d", id)
	return nil
}

// DeleteQuestion 删除问题
func (s *dailyQuestionService) DeleteQuestion(id uint) error {
	question, err := s.dailyQuestionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询问题失败, %w", err)
	}
	if question == nil {
		return bizerrors.New(bizerrors.CodeDailyQuestionNotFound, bizerrors.GetMessage(bizerrors.CodeDailyQuestionNotFound))
	}

	if err := s.dailyQuestionRepo.Delete(id); err != nil {
		return fmt.Errorf("删除问题失败, %w", err)
	}

	logger.Infof("删除问题成功, id: %d", id)
	return nil
}

// UpdateQuestionStatus 更新问题状态
func (s *dailyQuestionService) UpdateQuestionStatus(id uint, status int) error {
	if status < entity.DailyQuestionStatusDisabled || status > entity.DailyQuestionStatusScheduled {
		return bizerrors.New(bizerrors.CodeInvalidParams, "无效的问题状态")
	}
	question, err := s.dailyQuestionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询问题失败, %w", err)
	}
	if question == nil {
		return bizerrors.New(bizerrors.CodeDailyQuestionNotFound, bizerrors.GetMessage(bizerrors.CodeDailyQuestionNotFound))
	}

	question.Status = status
	if err := s.dailyQuestionRepo.Update(question); err != nil {
		return fmt.Errorf("更新问题状态失败, %w", err)
	}

	logger.Infof("更新问题状态成功, id: %d, status: %d", id, status)
	return nil
}

// toResponse 转换为响应DTO
func (s *dailyQuestionService) toResponse(q *entity.DailyQuestion) *response.DailyQuestionResponse {
	return &response.DailyQuestionResponse{
		ID:           q.ID,
		Question:     q.Question,
		Answer:       q.Answer,
		Date:         q.Date,
		LikeCount:    q.LikeCount,
		CommentCount: q.CommentCount,
		ViewCount:    q.ViewCount,
		Status:       q.Status,
		CreatedAt:    q.CreatedAt,
	}
}
