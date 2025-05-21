package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", GetEnv("REDIS_HOST", ""), GetEnv("REDIS_PORT", "")),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return err
	}

	RedisClient = client
	fmt.Println("Redis connection established successfully.")

	return nil
}
