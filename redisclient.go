package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisService() *redis.Client {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}

	fmt.Println("Connected to Redis:", pong)
	return redisClient
}
