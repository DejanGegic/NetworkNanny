package db

import (
	"os"
	"strings"

	"example.com/util"
)

func InitDB() DbInterface {
	util.LoadEnv()
	envDB := strings.ToLower(os.Getenv("DB_TYPE"))
	if envDB != "redis" {
		return &BadgerInstance{initBadger()}
	}
	RedisInstance := ConnectToRedisClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"), 0)
	return &RedisInstance
}
