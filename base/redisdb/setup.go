package redisdb

import (
	"app/base"
	"app/base/utils"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
)

func Configure() {
	port, err := strconv.Atoi(utils.Getenv("REDIS_PORT", "FILL"))
	if err != nil {
		panic(err)
	}
	host := utils.Getenv("REDIS_HOST", "FILL")
	utils.Log("host", host, "port", port).Info("Connecting to Redis")
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	check()
}

func check() {
	status := Rdb.Ping(base.Context)
	if status.Err() != nil {
		panic(status.Err())
	}
}
