package article

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

// Controller 文章控制器
type Controller struct {
	articleSvc service.ArticleService
}

// NewController 创建文章控制器
func NewController(articleSvc service.ArticleService) *Controller {
	return &Controller{articleSvc: articleSvc}
}

// GetArticleList 获取文章列表（前台）
func (c *Controller) GetArticleList(ctx *gin.Context) {
	var req request.ArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	result, err := c.articleSvc.GetArticleList(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取文章列表业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取文章列表失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetArticleDetail 获取文章详情（前台）
func (c *Controller) GetArticleDetail(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.BadRequest(ctx, "文章slug不能为空")
		return
	}

	result, err := c.articleSvc.GetArticleDetail(slug)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取文章详情业务错误", zap.String("slug", slug), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取文章详情失败", zap.String("slug", slug), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetArchives 获取文章归档
func (c *Controller) GetArchives(ctx *gin.Context) {
	result, err := c.articleSvc.GetArticleArchives()
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取文章归档业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取文章归档失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetAdminArticleList 获取文章列表（后台）
func (c *Controller) GetAdminArticleList(ctx *gin.Context) {
	var req request.ArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	result, err := c.articleSvc.GetAdminArticleList(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取后台文章列表业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取后台文章列表失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// GetAdminArticleDetail 获取文章详情（后台）
func (c *Controller) GetAdminArticleDetail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的文章ID")
		return
	}

	result, err := c.articleSvc.GetAdminArticleDetail(uint(id))
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("获取后台文章详情业务错误", zap.Uint64("id", id), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("获取后台文章详情失败", zap.Uint64("id", id), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, result)
}

// CreateArticle 创建文章
func (c *Controller) CreateArticle(ctx *gin.Context) {
	var req request.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	id, err := c.articleSvc.CreateArticle(&req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("创建文章业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("创建文章失败", zap.String("title", req.Title), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, gin.H{"id": id})
	ctx.Set("audit_created_id", uint(id))
}

// UpdateArticle 更新文章
func (c *Controller) UpdateArticle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的文章ID")
		return
	}

	var req request.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err = c.articleSvc.UpdateArticle(uint(id), &req)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("更新文章业务错误", zap.Uint64("id", id), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("更新文章失败", zap.Uint64("id", id), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// UpdateArticleStatus 更新文章状态（列表页快捷切换）
func (c *Controller) UpdateArticleStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的文章ID")
		return
	}

	var req request.UpdateArticleStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err = c.articleSvc.UpdateArticleStatus(uint(id), req.Status)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("更新文章状态业务错误", zap.Uint64("id", id), zap.String("status", req.Status), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("更新文章状态失败", zap.Uint64("id", id), zap.String("status", req.Status), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// DeleteArticle 删除文章
func (c *Controller) DeleteArticle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		response.BadRequest(ctx, "无效的文章ID")
		return
	}

	err = c.articleSvc.DeleteArticle(uint(id))
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("删除文章业务错误", zap.Uint64("id", id), zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("删除文章失败", zap.Uint64("id", id), zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// BatchDeleteArticles 批量删除文章
func (c *Controller) BatchDeleteArticles(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误")
		return
	}

	err := c.articleSvc.BatchDeleteArticles(req.IDs)
	if err != nil {
		if bizerrors.IsBizError(err) {
			logger.Warn("批量删除文章业务错误", zap.Error(err))
			response.BizError(ctx, err)
		} else {
			logger.Error("批量删除文章失败", zap.Error(err))
			response.ServerError(ctx, "服务器内部错误")
		}
		return
	}

	response.Success(ctx, nil)
}

// Search 搜索文章（前台，标题、内容、摘要模糊搜索）
func (c *Controller) Search(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		response.BadRequest(ctx, "搜索关键词不能为空")
		return
	}

	var req request.PageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 使用默认分页参数
		req.Page = 1
		req.PageSize = 10
	}

	result, err := c.articleSvc.Search(keyword, req.Page, req.GetPageSize())
	if err != nil {
		logger.Error("搜索文章失败", zap.String("keyword", keyword), zap.Error(err))
		response.ServerError(ctx, "服务器内部错误")
		return
	}

	response.Success(ctx, result)
}
