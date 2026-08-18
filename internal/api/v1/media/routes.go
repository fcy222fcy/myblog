package media

import (
	"blog/internal/middleware"
	blogjwt "blog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册媒体模块路由
func RegisterRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT) {
	// 需要登录的路由
	protected := rg.Group("")
	protected.Use(middleware.Auth(jwtInstance))
	{
		media := protected.Group("/media")
		{
			media.POST("/upload", controller.Upload)
			media.GET("", controller.List)
			media.DELETE("/:filename", controller.Delete)
		}
	}
}
