package tag

import (
	"blog/internal/middleware"
	blogjwt "blog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册标签模块路由
func RegisterRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT, bloggerUserID uint) {
	// 公开路由（无需登录）
	registerPublicRoutes(rg, controller)

	// 需要登录的路由
	registerProtectedRoutes(rg, controller, jwtInstance, bloggerUserID)
}

// registerPublicRoutes 注册公开路由
func registerPublicRoutes(rg *gin.RouterGroup, controller *Controller) {
	tags := rg.Group("/tags")
	{
		tags.GET("", controller.GetTagList)
	}
}

// registerProtectedRoutes 注册需要登录的路由
func registerProtectedRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT, bloggerUserID uint) {
	protected := rg.Group("")
	protected.Use(middleware.Auth(jwtInstance))
	protected.Use(middleware.RequireBlogger(bloggerUserID))
	{
		// 后台管理路由
		registerAdminRoutes(protected, controller)
	}
}

// registerAdminRoutes 注册后台管理路由
func registerAdminRoutes(rg *gin.RouterGroup, controller *Controller) {
	admin := rg.Group("/admin/tags")
	{
		admin.POST("", controller.CreateTag)
		admin.PUT("/:id", controller.UpdateTag)
		admin.DELETE("/:id", controller.DeleteTag)
	}
}
