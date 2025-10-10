package redis

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/kenziehh/cashflow-be/config"
)

func InitRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})

	log.Println("Redis connected successfully")
	return client
}