package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func SetUpRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	err = rdb.FlushDB(ctx).Err() // Use FlushAll(ctx) to clear all databases
	if err != nil {
		log.Fatalf("Failed to flush Redis DB: %v", err)
	}
	log.Println("Redis store cleared on startup")

	return rdb
}
