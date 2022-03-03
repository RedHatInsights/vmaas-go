package core

import (
	"app/manager/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter(handler gin.HandlerFunc) *gin.Engine {
	return InitRouterWithPath(handler, "/")
}

func InitRouterWithParams(handler gin.HandlerFunc, method, path string) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.RequestResponseLogger())
	router.Handle(method, path, handler)
	return router
}

func InitRouterWithPath(handler gin.HandlerFunc, path string) *gin.Engine {
	return InitRouterWithParams(handler, "GET", path)
}
