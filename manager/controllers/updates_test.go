package controllers

import (
	"app/base"
	"app/base/core"
	"app/base/redisdb"
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

func addTestingRedisData() {
	redisdb.Configure()
	redisdb.Rdb.FlushAll(base.Context)
	redisdb.Rdb.SAdd(base.Context, "u:pA", 1, 2, 3, 4, 5, 6)
	redisdb.Rdb.SAdd(base.Context, "a:x86", 2, 4, 6)
	redisdb.Rdb.SAdd(base.Context, "a:i686", 1, 3, 5)
	redisdb.Rdb.Set(base.Context, "1", "pA-1-1.el7.i686 ER7-1", 0)
	redisdb.Rdb.Set(base.Context, "2", "pA-1-1.el7.x86 ER7-1", 0)
	redisdb.Rdb.Set(base.Context, "3", "pA-1.1-1.el7.i686 ER7-2", 0)
	redisdb.Rdb.Set(base.Context, "4", "pA-1.1-1.el7.x86 ER7-2", 0)
	redisdb.Rdb.Set(base.Context, "5", "pA-2-1.el8.i686 ER8-1", 0)
	redisdb.Rdb.Set(base.Context, "6", "pA-2-1.el8.x86 ER8-1", 0)
	redisdb.Rdb.SAdd(base.Context, "r:rhel7", 1, 2, 3, 4)
	redisdb.Rdb.SAdd(base.Context, "r:rhel8", 5, 6)
}

func TestUpdatesDefault(t *testing.T) {
	core.SetupTestEnvironment()

	addTestingRedisData()
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

	addTestingRedisData()
	data := `{"package_list": "invalid value"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(UpdatesHandler, "POST", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "Invalid request body")
}
