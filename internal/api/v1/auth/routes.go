package auth

import (
	"blog/internal/middleware"
	blogjwt "blog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册认证模块路由
func RegisterRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Register)
	}

	// 需要登录的路由
	protected := rg.Group("")
	protected.Use(middleware.Auth(jwtInstance))
	{
		auth := protected.Group("/auth")
		{
			auth.PUT("/password", controller.ChangePassword)
		}
	}
}
