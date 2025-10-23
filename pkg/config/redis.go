package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	// Gunakan env untuk mode
	mode := os.Getenv("APP_MODE") // contoh: dev / prod

	opt := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		if mode == "dev" {
			log.Println("⚠️ Redis not running — skipping cache for dev mode")
			RedisClient = nil
			return
		} else {
			log.Fatalf("❌ Redis connection failed: %v", err)
		}
	}

	log.Println("✅ Redis connected successfully")
	RedisClient = client
}
