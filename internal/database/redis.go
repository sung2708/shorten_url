package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func InitRedis(host string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("Failed to connect to redis")
	}
	return rdb
}
