package comment

import (
	"blog/internal/middleware"
	blogjwt "blog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册评论模块路由
func RegisterRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT) {
	// 公开路由（无需登录）
	registerPublicRoutes(rg, controller, jwtInstance)

	// 需要登录的路由
	registerProtectedRoutes(rg, controller, jwtInstance)
}

// registerPublicRoutes 注册公开路由
func registerPublicRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT) {
	comments := rg.Group("/comments")
	comments.Use(middleware.OptionalAuth(jwtInstance))
	{
		comments.GET("/article/:articleId", controller.GetCommentsByArticle)
		comments.POST("", controller.CreateComment)
		comments.POST("/:id/like", controller.LikeComment)
		comments.DELETE("/:id/like", controller.UnlikeComment)
	}
}

// registerProtectedRoutes 注册需要登录的路由
func registerProtectedRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT) {
	protected := rg.Group("")
	protected.Use(middleware.Auth(jwtInstance))
	{
		// 后台管理路由
		registerAdminRoutes(protected, controller)
	}
}

// registerAdminRoutes 注册后台管理路由
func registerAdminRoutes(rg *gin.RouterGroup, controller *Controller) {
	admin := rg.Group("/admin/comments")
	{
		admin.GET("", controller.GetAdminCommentList)
		admin.PUT("/:id/status", controller.UpdateCommentStatus)
		admin.DELETE("/:id", controller.DeleteComment)
		admin.POST("/batch-delete", controller.BatchDeleteComments)
	}
}
