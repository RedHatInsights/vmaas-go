package core

import (
	"app/base/database"
	"app/base/redisdb"
	"app/base/utils"
)

func ConfigureApp() {
	utils.ConfigureLogging()
	database.Configure()
	redisdb.Configure()
}

func SetupTestEnvironment() {
	utils.SetDefaultEnvOrFail("LOG_LEVEL", "debug")
	ConfigureApp()
}
