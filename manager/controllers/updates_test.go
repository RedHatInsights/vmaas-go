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
	data := `{"package_list": ["pA-1-1.el7.x86", "unknown-1-1.el8.i686"],
              "repository_list": ["rhel7"]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(UpdatesHandler, "POST", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp UpdatesResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp.UpdateList))
	assert.Equal(t, 1, len(resp.UpdateList["pA-1-1.el7.x86"]))
	assert.Equal(t, []string{"pA-1.1-1.el7.x86", "ER7-2"}, resp.UpdateList["pA-1-1.el7.x86"][0])
	assert.Equal(t, 0, len(resp.UpdateList["unknown-1-1.el8.i686"]))
}

func TestUpdatesUnknownRepo(t *testing.T) {
	core.SetupTestEnvironment()

	addTestingRedisData()
	data := `{"package_list": ["pA-1-1.el7.i686", "unknown-1-1.el8.i686"],
              "repository_list": ["unknown"]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(UpdatesHandler, "POST", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp UpdatesResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp.UpdateList))
	assert.Equal(t, 0, len(resp.UpdateList["pA-1-1.el7.i686"]))
	assert.Equal(t, 0, len(resp.UpdateList["unknown-1-1.el8.i686"]))
}

func TestUpdatesNoRepo(t *testing.T) {
	core.SetupTestEnvironment()

	addTestingRedisData()
	data := `{"package_list": ["pA-1-1.el7.i686", "unknown-1-1.el8.i686"]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(UpdatesHandler, "POST", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp UpdatesResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp.UpdateList))
	assert.Equal(t, 2, len(resp.UpdateList["pA-1-1.el7.i686"]))
	assert.Equal(t, []string{"pA-1.1-1.el7.i686", "ER7-2"}, resp.UpdateList["pA-1-1.el7.i686"][0])
	assert.Equal(t, []string{"pA-2-1.el8.i686", "ER8-1"}, resp.UpdateList["pA-1-1.el7.i686"][1])
	assert.Equal(t, 0, len(resp.UpdateList["unknown-1-1.el8.i686"]))
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

// TestGetReposKey

func TestGetReposKey(t *testing.T) {
	core.SetupTestEnvironment()
	addTestingRedisData()

	keyArr, err := getReposKey([]string{"rhel7", "rhel8"})
	assert.Nil(t, err)
	assert.Equal(t, "r:rhel7+r:rhel8", keyArr[0])
	result := redisdb.Rdb.SMembers(base.Context, keyArr[0])
	indices, err := result.Result()
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, indices)

	keyArr, err = getReposKey([]string{"rhel7"})
	assert.Nil(t, err)
	assert.Equal(t, "r:rhel7", keyArr[0])

	keyArr, err = getReposKey([]string{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(keyArr))
}

// TestGetNevraUpdates

func TestGetNevraUpdates(t *testing.T) {
	core.SetupTestEnvironment()
	addTestingRedisData()

	nevraUpdatesIDs, nevra := getNevraUpdates("pA-1-1.el7.x86", nil)
	assert.Equal(t, []string{"2", "4", "6"}, nevraUpdatesIDs)
	assert.Equal(t, "pA", nevra.Name)
}

// TestGetNevraErratumUpdates

func TestGetNevraErratumUpdates(t *testing.T) {
	core.SetupTestEnvironment()
	addTestingRedisData()

	packageList := []string{"pA-1-1.el7.x86"}
	repositoryList := []string{"rhel7"}
	nevraErratumUpdates, nevras, counts, err := getNevraErratumUpdates(packageList, repositoryList)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(nevraErratumUpdates))
	assert.Equal(t, 1, len(nevras))
	assert.Equal(t, []int{2}, counts)
}

func TestGetNevraErratumUpdatesNoRepo(t *testing.T) {
	core.SetupTestEnvironment()
	addTestingRedisData()

	packageList := []string{"pA-1-1.el7.x86"}
	nevraErratumUpdates, nevras, counts, err := getNevraErratumUpdates(packageList, nil)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(nevraErratumUpdates))
	assert.Equal(t, 1, len(nevras))
	assert.Equal(t, []int{3}, counts)
}

func TestGetNevraErratumUpdatesTwoPackages(t *testing.T) {
	core.SetupTestEnvironment()
	addTestingRedisData()

	packageList := []string{"pA-1-1.el7.x86", "pA-1-1.el7.x86"}
	repositoryList := []string{"rhel7"}
	nevraErratumUpdates, nevras, counts, err := getNevraErratumUpdates(packageList, repositoryList)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(nevraErratumUpdates))
	assert.Equal(t, 2, len(nevras))
	assert.Equal(t, []int{2, 2}, counts)
}

func TestGetNevraErratumUpdatesUnknownRepo(t *testing.T) {
	core.SetupTestEnvironment()
	addTestingRedisData()

	packageList := []string{"pA-1-1.el7.x86", "pB-1-1.el7.x86"}
	repositoryList := []string{"unknown-repo"}
	nevraErratumUpdates, nevras, counts, err := getNevraErratumUpdates(packageList, repositoryList)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(nevraErratumUpdates))
	assert.Equal(t, 2, len(nevras))
	assert.Equal(t, []int{0, 0}, counts)
}

// FilterNevraUpdates

func TestFilterNevraUpdates(t *testing.T) {
	installedEvra := "1-1.el7.x86"
	var nevraErratumUpdates []interface{}
	nevraErratumUpdates = append(nevraErratumUpdates, "pA-0-1.el7.x86 ER0", "pA-1-1.el7.x86 ER1", "pA-2-1.el7.x86 ER2", "pA-3-1.el7.x86 ER3")
	nevraUpdates := filterNevraUpdates(installedEvra, nevraErratumUpdates)
	assert.Equal(t, 2, len(nevraUpdates))
	assert.Equal(t, []string{"pA-2-1.el7.x86", "ER2"}, nevraUpdates[0])
	assert.Equal(t, []string{"pA-3-1.el7.x86", "ER3"}, nevraUpdates[1])
}
