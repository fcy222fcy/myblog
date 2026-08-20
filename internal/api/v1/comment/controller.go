package comment

import (
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	bizerrors "blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Controller 评论控制器
type Controller struct {
	commentSvc service.CommentService
}

// NewController 创建评论控制器
func NewController(commentSvc service.CommentService) *Controller {
	return &Controller{commentSvc: commentSvc}
}

// GetCommentsByArticle 获取文章评论列表
func (c *Controller) GetCommentsByArticle(ctx *gin.Context) {
	articleParam := ctx.Param("articleId")
	logger.Infof("[评论] 收到请求, articleParam=%s, query=%s", articleParam, ctx.Request.URL.RawQuery)

	var req request.CommentListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Warnf("[评论] 参数绑定失败: %v", err)
		response.BadRequest(ctx, "参数错误")
		return
	}

	result, err := c.commentSvc.GetCommentsByArticle(articleParam, &req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取评论列表业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取评论列表失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// CreateComment 创建评论
func (c *Controller) CreateComment(ctx *gin.Context) {
	var req request.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("[评论] 创建评论参数绑定失败: %v", err)
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	logger.Infof("[评论] 创建评论, nickname=%s, articleID=%d", req.Nickname, req.ArticleID)

	userID := middleware.GetUserID(ctx)
	ip := ctx.ClientIP()
	ua := ctx.Request.UserAgent()

	id, err := c.commentSvc.CreateComment(&req, userID, ip, ua)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("创建评论业务错误", zap.String("nickname", req.Nickname), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("创建评论失败", zap.String("nickname", req.Nickname), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, gin.H{"id": id})
	ctx.Set("audit_created_id", uint(id))
}

// GetAdminCommentList 获取评论列表（后台）
func (c *Controller) GetAdminCommentList(ctx *gin.Context) {
	var req request.CommentListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	result, err := c.commentSvc.GetAdminCommentList(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取后台评论列表业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取后台评论列表失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// UpdateCommentStatus 更新评论状态
func (c *Controller) UpdateCommentStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的评论ID")
		return
	}

	var req request.UpdateCommentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err = c.commentSvc.UpdateCommentStatus(uint(id), req.Status)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("更新评论状态业务错误", zap.Uint64("id", id), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("更新评论状态失败", zap.Uint64("id", id), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// DeleteComment 删除评论
func (c *Controller) DeleteComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的评论ID")
		return
	}

	err = c.commentSvc.DeleteComment(uint(id))
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("删除评论业务错误", zap.Uint64("id", id), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("删除评论失败", zap.Uint64("id", id), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// BatchDeleteComments 批量删除评论
func (c *Controller) BatchDeleteComments(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err := c.commentSvc.BatchDeleteComments(req.IDs)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("批量删除评论业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("批量删除评论失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// LikeComment 点赞评论
func (c *Controller) LikeComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的评论ID")
		return
	}

	visitorIP := ctx.ClientIP()

	err = c.commentSvc.LikeComment(uint(id), visitorIP)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("点赞评论业务错误", zap.Uint64("id", id), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("点赞评论失败", zap.Uint64("id", id), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// UnlikeComment 取消点赞评论
func (c *Controller) UnlikeComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的评论ID")
		return
	}

	visitorIP := ctx.ClientIP()

	err = c.commentSvc.UnlikeComment(uint(id), visitorIP)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("取消点赞评论业务错误", zap.Uint64("id", id), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("取消点赞评论失败", zap.Uint64("id", id), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}
