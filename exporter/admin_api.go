package exporter

import (
	"app/base"
	"app/base/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func runAdminAPI() {
	app := gin.New()

	app.GET("/sync", sync)

	err := utils.RunServer(base.Context, app, 9999)

	if err != nil {
		utils.Log("err", err.Error()).Error()
		panic(err)
	}
}

func sync(c *gin.Context) {
	utils.Log().Info("manual syncing called...")
	err := syncData()
	if err != nil {
		utils.Log("err", err.Error()).Error("manual called syncing failed")
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	utils.Log().Info("manual syncing finished successfully")
	c.JSON(http.StatusOK, "OK")
}
