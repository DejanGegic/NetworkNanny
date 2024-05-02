package db

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisInstance struct {
	*redis.Client
}

func ConnectToRedisClient(host string, port string, password string, db int) RedisInstance {
	log.Println("Connecting to redis at", host)
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	instance := RedisInstance{
		client,
	}
	if err := instance.Ping().Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Redis connected")
	return instance
}

func (r RedisInstance) SetEvictionPolicy(policy string) error {
	err := r.ConfigSet("maxmemory-policy", policy).Err()
	if err != nil {
		return err
	}
	return err
}

func (r RedisInstance) IncrementValue(key string) error {
	err := r.Incr(key).Err()
	return err
}

// RedisInstance implements DbInterface
// badger should also have a struct that implements DbInterface
func (r RedisInstance) Write(key string, value string) error {
	err := r.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return err
}

func (r RedisInstance) WriteTTL(key string, value string, tll time.Duration) error {
	return r.Set(key, value, tll).Err()
}

func (r RedisInstance) Read(key string) (string, error) {

	val, err := r.Get(key).Result()
	if err != nil {
		val = "0"
	}
	return val, err
}

func (r RedisInstance) ReadTTL(key string) (string, time.Duration, error) {
	val, err := r.Get(key).Result()
	if err != nil {
		val = ""
		return val, 0, err
	}
	// convert the string to an int
	ttl, err := r.TTL(key).Result()
	return val, ttl, err
}

func (r RedisInstance) CheckHealth() error {
	return r.Ping().Err()
}
