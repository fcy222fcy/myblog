package auth

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	bizerrors "blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Controller 认证控制器
type Controller struct {
	authSvc service.AuthService
}

// NewController 创建认证控制器
func NewController(authSvc service.AuthService) *Controller {
	return &Controller{authSvc: authSvc}
}

// Login 用户登录
func (c *Controller) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	// 自定义校验：邮箱和用户名至少填一个
	if !req.ValidateLogin() {
		response.BadRequest(ctx, "请输入邮箱或用户名")
		return
	}

	result, err := c.authSvc.Login(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("登录失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// Register 用户注册
func (c *Controller) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err := c.authSvc.Register(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("注册失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// ChangePassword 修改密码
func (c *Controller) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未登录")
		return
	}

	var req request.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	err := c.authSvc.ChangePassword(userID.(uint), &req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("修改密码业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("修改密码失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}
