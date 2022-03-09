package exporter

import (
	"app/base"
	"app/base/core"
	"app/base/utils"
	"os"
	"time"
)

var (
	enableSyncOnStart bool
)

func configure() {
	core.ConfigureApp()
	enableSyncOnStart = utils.GetBoolEnvOrDefault("ENABLE_SYNC_ON_START", false)
}

func handleContextCancel(fn func()) {
	go func() {
		<-base.Context.Done()
		utils.Log().Info("stopping vmaas_sync")
		fn()
	}()
}

func waitAndExit() {
	time.Sleep(time.Second) // give some time to close eventual db connections
	os.Exit(0)
}

func syncOnStartIfSet() {
	if enableSyncOnStart {
		err := syncData()
		if err != nil {
			utils.Log("err", err.Error()).Error("unable to sync data on start")
		}
	}
}

func RunVmaasSync() {
	handleContextCancel(waitAndExit)
	configure()

	syncOnStartIfSet() // sync data start if configured

	runAdminAPI()
}

func syncData() error {
	utils.Log().Info("Data sync started")

	err := syncPackages()
	if err != nil {
		return err
	}

	err = syncRepos()
	if err != nil {
		return err
	}

	utils.Log().Info("Data sync finished successfully")
	return nil
}
