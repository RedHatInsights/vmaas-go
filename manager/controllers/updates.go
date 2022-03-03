package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdatesRequest struct {
	PackageList    []string `json:"package_list"`
	RepositoryList []string `json:"repository_list"`
}

type UpdatesResponse struct {
	UpdateList map[string][][]string `json:"update_list"`
}

// @Summary Show updates
// @Description Show updates
// @ID listAdvisories
// @Security RhIdentity
// @Accept   json
// @Produce  json
// @Success 200 {object} UpdatesResponse
// @Router /api/patch/v1/updates [post]
func UpdatesHandler(c *gin.Context) {
	var request UpdatesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		LogAndRespBadRequest(c, err, "Invalid request body: "+err.Error())
		return
	}

	// Fake response
	var resp = UpdatesResponse{
		UpdateList: map[string][][]string{
			"pA-1-1.el7.i686":      {{"pA-1.1-1.el7.i686", "ER7-2"}},
			"unknown-1-1.el8.i686": {}},
	}
	c.JSON(http.StatusOK, &resp)
}
