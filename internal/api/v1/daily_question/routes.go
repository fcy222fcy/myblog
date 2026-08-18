package daily_question

import (
	"blog/internal/middleware"
	blogjwt "blog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册每日一问模块路由
func RegisterRoutes(rg *gin.RouterGroup, controller *Controller, jwtInstance *blogjwt.JWT) {
	// 公开路由（无需登录）
	registerPublicRoutes(rg, controller)

	// 需要登录的路由
	registerProtectedRoutes(rg, controller, jwtInstance)
}

// registerPublicRoutes 注册公开路由
func registerPublicRoutes(rg *gin.RouterGroup, controller *Controller) {
	dailyQ := rg.Group("/daily-questions")
	{
		dailyQ.GET("/latest", controller.GetLatestQuestion)
		dailyQ.GET("/all", controller.GetAllPublishedQuestions)
		dailyQ.GET("/date/:date", controller.GetQuestionByDate)
		dailyQ.GET("/previous/:date", controller.GetPreviousQuestion)
		dailyQ.GET("/next/:date", controller.GetNextQuestion)
		dailyQ.POST("/:id/like", controller.LikeQuestion)
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
	admin := rg.Group("/admin/daily-questions")
	{
		admin.GET("", controller.GetAdminQuestionList)
		admin.POST("", controller.CreateQuestion)
		admin.PUT("/:id", controller.UpdateQuestion)
		admin.DELETE("/:id", controller.DeleteQuestion)
		admin.PUT("/:id/status", controller.UpdateQuestionStatus)
	}
}
