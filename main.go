package main

import (
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/config"
	"github.com/RedHatInsights/vmaas-go/app/database"
	"github.com/RedHatInsights/vmaas-go/app/utils"
	"github.com/RedHatInsights/vmaas-go/app/webserver"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	config.SQLiteFilePath = os.Args[1]
	database.Configure()
	cache.C = cache.LoadCache()
	cache.C.Inspect()

	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	utils.PrintMemUsage()
	webserver.Run()
}
