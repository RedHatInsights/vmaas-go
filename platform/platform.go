package platform

import (
	"app/base"
	"app/base/utils"
	"app/manager/middlewares"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var websockets []chan string

func addWebsocket() chan string {
	ws := make(chan string, 100)
	websockets = append(websockets, ws)
	return ws
}

func platformMock() {
	utils.Log().Info("Platform mock starting")
	app := gin.New()
	app.Use(middlewares.RequestResponseLogger())
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	initVMaaS(app)

	// Control endpoint handler
	app.POST("/control/sync", mockSyncHandler)

	err := utils.RunServer(base.Context, app, 9001)
	if err != nil {
		panic(err)
	}
}

func mockSyncHandler(_ *gin.Context) {
	utils.Log().Info("Mocking VMaaS sync event")
	// Force connected websocket clients to refresh
	for _, ws := range websockets {
		ws <- "sync"
	}
}

func RunPlatformMock() {
	platformMock()
}
