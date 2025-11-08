package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func InitRedis(host string, user string, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Username: user,
		Password: password,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	return rdb
}
