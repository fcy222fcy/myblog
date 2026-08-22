package content_view

import (
	"blog/internal/repository"
	"blog/internal/service"
	"blog/pkg/response"
	"blog/pkg/visitor"
	"errors"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service service.ContentViewService
	secret  string
}

func NewController(viewService service.ContentViewService, secret string) *Controller {
	return &Controller{service: viewService, secret: secret}
}

func (c *Controller) Record(ctx *gin.Context) {
	var req struct {
		ContentType string `json:"content_type" binding:"required,oneof=article daily_question"`
		ContentID   uint   `json:"content_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || req.ContentID == 0 {
		response.BadRequest(ctx, "参数错误")
		return
	}

	visitorKey, err := visitor.DeriveKey(c.secret, ctx.GetHeader("X-Visitor-ID"), ctx.ClientIP())
	if err != nil {
		response.BadRequest(ctx, "无法识别访客")
		return
	}

	result, err := c.service.Record(req.ContentType, req.ContentID, visitorKey)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUnsupportedContentType):
			response.BadRequest(ctx, "参数错误")
		case errors.Is(err, repository.ErrViewTargetNotFound):
			response.NotFound(ctx, "内容不存在")
		default:
			response.ServerError(ctx, "浏览量记录失败")
		}
		return
	}

	response.Success(ctx, result)
}
