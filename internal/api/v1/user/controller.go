package user

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	bizerrors "blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Controller 用户控制器
type Controller struct {
	userSvc service.UserService
}

// NewController 创建用户控制器
func NewController(userSvc service.UserService) *Controller {
	return &Controller{userSvc: userSvc}
}

// GetPublicUserInfo 获取公开用户信息（无需登录）
func (c *Controller) GetPublicUserInfo(ctx *gin.Context) {
	// 获取第一个管理员用户（博客主人）
	user, err := c.userSvc.GetFirstAdmin()
	if err != nil || user == nil {
		// 返回默认信息
		response.Success(ctx, gin.H{
			"nickname": "Fu Chengyan",
			"bio":      "日常落灰的个人博客，擅长面向搜索引擎编程。分享 Golang 开发、AI 和 NAS 折腾经验",
			"avatar":   "",
		})
		return
	}

	response.Success(ctx, gin.H{
		"nickname": user.Nickname,
		"bio":      user.Bio,
		"avatar":   user.Avatar,
	})
}

// GetUserInfo 获取用户信息
func (c *Controller) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未登录")
		return
	}

	user, err := c.userSvc.GetUserByID(userID.(uint))
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取用户信息业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取用户信息失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, user)
}

// UpdateUserInfo 更新用户信息
func (c *Controller) UpdateUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未登录")
		return
	}

	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err := c.userSvc.UpdateUser(userID.(uint), &req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("更新用户信息业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("更新用户信息失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}
