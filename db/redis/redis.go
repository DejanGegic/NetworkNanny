package redis

import (
	"example.com/l"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //
		DB:       0,  // use default DB

	})

	pong, err := rdb.Ping().Result()
	if err != nil {
		l.Error(err)
		panic(err)
	}
	l.Info(pong)
}
