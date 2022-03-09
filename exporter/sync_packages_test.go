package exporter

import (
	"app/base"
	"app/base/core"
	"app/base/redisdb"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncPackages(t *testing.T) {
	core.SetupTestEnvironment()

	redisdb.Rdb.FlushAll(base.Context)
	err := syncPackages()
	assert.Nil(t, err)

	checkPackageDataInRedis(t)
}

func checkPackageDataInRedis(t *testing.T) {
	nevras, err := redisdb.Rdb.MGet(base.Context, "301", "302", "303", "304", "305", "306", "307").Result()
	assert.Nil(t, err)
	assert.Equal(t, 7, len(nevras))
	assert.Equal(t, "pkg-sec-errata1-1:1-1.noarch errata3", nevras[0])
	assert.Equal(t, "pkg-sec-errata1-1:1-2.noarch errata1", nevras[1])
	assert.Equal(t, "pkg-sec-errata1-1:1-3.noarch -", nevras[2])
	assert.Equal(t, "pkg-no-sec-errata2-2:2-2.noarch errata2", nevras[3])
	assert.Equal(t, "pkg-errata-cve3-3:3-3.noarch errata3", nevras[4])
	assert.Equal(t, "pkg-errata-cve3-4:4-4.noarch errata3", nevras[5])
	assert.Equal(t, "pkg-sec-errata4-2:2-2.noarch errata1", nevras[6])

	nameIDs, err := redisdb.Rdb.SMembers(base.Context, "u:pkg-sec-errata1").Result()
	assert.Nil(t, err)
	sort.Strings(nameIDs)
	assert.Equal(t, []string{"301", "302", "303"}, nameIDs)

	archIDs, err := redisdb.Rdb.SMembers(base.Context, "a:noarch").Result()
	assert.Nil(t, err)
	sort.Strings(archIDs)
	assert.Equal(t, []string{"301", "302", "303", "304", "305", "306", "307"}, archIDs)
}
