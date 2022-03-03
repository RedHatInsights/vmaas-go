package core

import (
	"app/base/database"
	"app/base/utils"
)

func ConfigureApp() {
	utils.ConfigureLogging()
	database.Configure()
}

func SetupTestEnvironment() {
	utils.SetDefaultEnvOrFail("LOG_LEVEL", "debug")
	ConfigureApp()
}
