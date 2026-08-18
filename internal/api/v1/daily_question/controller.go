package daily_question

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	bizerrors "blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Controller 每日一问控制器
type Controller struct {
	dailyQSvc service.DailyQuestionService
}

// NewController 创建每日一问控制器
func NewController(dailyQSvc service.DailyQuestionService) *Controller {
	return &Controller{dailyQSvc: dailyQSvc}
}

// GetLatestQuestion 获取最新问题
func (c *Controller) GetLatestQuestion(ctx *gin.Context) {
	result, err := c.dailyQSvc.GetLatestQuestion()
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取最新问题业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取最新问题失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetAllPublishedQuestions 获取所有已发布问题列表
func (c *Controller) GetAllPublishedQuestions(ctx *gin.Context) {
	result, err := c.dailyQSvc.GetAllPublishedQuestions()
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取已发布问题列表业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取已发布问题列表失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetQuestionByDate 根据日期获取问题
func (c *Controller) GetQuestionByDate(ctx *gin.Context) {
	date := ctx.Param("date")
	if date == "" {
		response.BadRequest(ctx, "日期不能为空")
		return
	}

	result, err := c.dailyQSvc.GetQuestionByDate(date)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("根据日期获取问题业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("根据日期获取问题失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetPreviousQuestion 获取前一天的问题
func (c *Controller) GetPreviousQuestion(ctx *gin.Context) {
	date := ctx.Param("date")
	if date == "" {
		response.BadRequest(ctx, "日期不能为空")
		return
	}

	result, err := c.dailyQSvc.GetPreviousQuestion(date)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取前一天问题业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取前一天问题失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetNextQuestion 获取后一天的问题
func (c *Controller) GetNextQuestion(ctx *gin.Context) {
	date := ctx.Param("date")
	if date == "" {
		response.BadRequest(ctx, "日期不能为空")
		return
	}

	result, err := c.dailyQSvc.GetNextQuestion(date)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取后一天问题业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取后一天问题失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// LikeQuestion 问题点赞
func (c *Controller) LikeQuestion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的问题ID")
		return
	}

	likeCount, err := c.dailyQSvc.LikeQuestion(uint(id), ctx.ClientIP())
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("问题点赞业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("问题点赞失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, gin.H{"like_count": likeCount})
}

// GetAdminQuestionList 获取问题列表（后台）
func (c *Controller) GetAdminQuestionList(ctx *gin.Context) {
	var req request.DailyQuestionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	result, err := c.dailyQSvc.GetAdminQuestionList(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取后台问题列表业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取后台问题列表失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// CreateQuestion 创建问题
func (c *Controller) CreateQuestion(ctx *gin.Context) {
	var req request.CreateDailyQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	id, err := c.dailyQSvc.CreateQuestion(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("创建问题业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("创建问题失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, gin.H{"id": id})
	ctx.Set("audit_created_id", uint(id))
}

// UpdateQuestion 更新问题
func (c *Controller) UpdateQuestion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的问题ID")
		return
	}

	var req request.UpdateDailyQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err = c.dailyQSvc.UpdateQuestion(uint(id), &req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("更新问题业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("更新问题失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// DeleteQuestion 删除问题
func (c *Controller) DeleteQuestion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的问题ID")
		return
	}

	err = c.dailyQSvc.DeleteQuestion(uint(id))
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("删除问题业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("删除问题失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// UpdateQuestionStatus 更新问题状态
func (c *Controller) UpdateQuestionStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的问题ID")
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}
	if req.Status != 0 && req.Status != 1 {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err = c.dailyQSvc.UpdateQuestionStatus(uint(id), req.Status)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("更新问题状态业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("更新问题状态失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}
