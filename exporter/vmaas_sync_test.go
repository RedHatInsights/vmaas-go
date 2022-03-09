package exporter

import (
	"app/base"
	"app/base/core"
	"app/base/redisdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncData(t *testing.T) {
	core.SetupTestEnvironment()

	redisdb.Rdb.FlushAll(base.Context)
	err := syncData()
	assert.Nil(t, err)

	checkPackageDataInRedis(t)
	checkReposDataInRedis(t)
}
