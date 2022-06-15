package config

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	config := conf.Redis
	if config == nil {
		panic("[FAILED] config file failed - Redis")
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:           config.Password,
		DB:                 config.DB,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: time.Second,
	})
	if RedisClient == nil {
		panic("[FAILED] init Redis failed")
	}
	log.Println("[INFO] init Redis finished successfully")
}
