package routes

import (
	"app/docs"
	"app/manager/controllers"

	"github.com/gin-gonic/gin"
)

func InitAPI(api *gin.RouterGroup, config docs.EndpointsConfig) {
	api.POST("/updates", controllers.UpdatesHandler)
}
