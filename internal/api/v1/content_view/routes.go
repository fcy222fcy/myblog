package content_view

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, controller *Controller) {
	rg.POST("/views", controller.Record)
}
