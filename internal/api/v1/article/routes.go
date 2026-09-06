package article

import (
	"blog/internal/middleware"
	blogjwt "blog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文章模块路由
func RegisterRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT) {
	// 公开路由（无需登录）
	registerPublicRoutes(rg, controller)

	// 需要登录的路由
	registerProtectedRoutes(rg, controller, jwtInstance)
}

// registerPublicRoutes 注册公开路由
func registerPublicRoutes(rg *gin.RouterGroup, controller *Controller) {
	articles := rg.Group("/articles")
	{
		articles.GET("", controller.GetArticleList)
		articles.GET("/search", controller.Search)
		articles.GET("/archives", controller.GetArchives)
		articles.GET("/:slug", controller.GetArticleDetail)
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
	admin := rg.Group("/admin/articles")
	{
		admin.GET("", controller.GetAdminArticleList)
		admin.GET("/:id", controller.GetAdminArticleDetail)
		admin.POST("", controller.CreateArticle)
		admin.PUT("/:id", controller.UpdateArticle)
		admin.PUT("/:id/status", controller.UpdateArticleStatus)
		admin.DELETE("/:id", controller.DeleteArticle)
		admin.POST("/batch-delete", controller.BatchDeleteArticles)
	}
}
