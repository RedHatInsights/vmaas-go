package exporter

import (
	"app/base"
	"app/base/core"
	"app/base/redisdb"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncRepos(t *testing.T) {
	core.SetupTestEnvironment()

	redisdb.Rdb.FlushAll(base.Context)
	err := syncRepos()
	assert.Nil(t, err)

	checkReposDataInRedis(t)
}

func checkReposDataInRedis(t *testing.T) {
	ID1s, err := redisdb.Rdb.SMembers(base.Context, "r:content-set-name-1").Result()
	assert.Nil(t, err)
	sort.Strings(ID1s)
	assert.Equal(t, []string{"301", "302", "303", "304", "305", "306"}, ID1s)

	ID2s, err := redisdb.Rdb.SMembers(base.Context, "r:content-set-name-2").Result()
	assert.Nil(t, err)
	sort.Strings(ID2s)
	assert.Equal(t, []string{"306", "307"}, ID2s)
}
