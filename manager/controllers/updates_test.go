package controllers

import (
	"app/base/core"
	"app/base/utils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	utils.TestLoadEnv("conf/manager.env")
}

func TestUpdatesDefault(t *testing.T) {
	core.SetupTestEnvironment()
	data := `{"package_list": ["pA-1-1.el7.i686", "unknown-1-1.el8.i686"],
              "repository_list": ["rhel7"]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(UpdatesHandler, "POST", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp UpdatesResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp.UpdateList))
}

func TestUpdatesInvalidRequest(t *testing.T) {
	core.SetupTestEnvironment()
	data := `{"package_list": "invalid value"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(UpdatesHandler, "POST", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "Invalid request body")
}
