package database_admin //nolint:revive,stylecheck

import (
	"app/base/database"
	"app/base/utils"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setCmdAuth(cmd *exec.Cmd) {
	cmd.Args = append(cmd.Args,
		"-h", utils.GetenvOrFail("DB_HOST"),
		"-p", utils.GetenvOrFail("DB_PORT"),
		"-U", utils.GetenvOrFail("DB_USER"),
		"-d", utils.GetenvOrFail("DB_NAME"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", utils.GetenvOrFail("DB_PASSWD")))
}

func TestInitSchema(t *testing.T) {
	utils.SkipWithoutDB(t)
	database.Configure()

	err := database.ExecFile("./schema/clear_db.sql")
	assert.NoError(t, err)

	err = database.ExecFile("./schema/create_schema.sql")
	assert.NoError(t, err)

	dumpCmd := exec.Command("pg_dump", "-O")
	setCmdAuth(dumpCmd)

	fromScratch, err := dumpCmd.Output()
	assert.NoError(t, err)

	scratchLines := strings.SplitAfter(string(fromScratch), "\n")
	assert.LessOrEqual(t, 40, len(scratchLines))
}
