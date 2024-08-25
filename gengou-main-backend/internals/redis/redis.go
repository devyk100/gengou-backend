package redis

import (
	"github.com/devyk100/gengou-db/pkg/redis_internal"
	"os"
	"time"
)

var Instance *redis_internal.RedisInstance
var err error

func RedisInit() {
	dsn := os.Getenv("REDIS_URL")
	Instance, err = redis_internal.Init(dsn, time.Second*100)
	if err != nil {
		return
	}
}

func RedisClose() {
	err := Instance.Close()
	if err != nil {
		panic(err.Error())
		return
	}
}
