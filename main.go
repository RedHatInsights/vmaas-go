package main

import (
	"app/base"
	"app/base/utils"
	"app/manager"
	"app/platform"
	"log"
	"os"
)

func main() {
	base.HandleSignals()

	defer utils.LogPanics(true)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "manager":
			manager.RunManager()
			return
		case "platform":
			platform.RunPlatformMock()
			return
		case "print_clowder_params":
			utils.PrintClowderParams()
			return
		}
	}
	log.Panic("You need to provide a command")
}
