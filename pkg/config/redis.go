package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	Ctx        = context.Background()
	RedisClient *redis.Client
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // default Redis port
		Password: "",               // kosong jika tidak pakai password
		DB:       0,                // gunakan DB 0
	})

	// Test koneksi
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ Gagal konek Redis: %v", err))
	}

	fmt.Println("✅ Redis connected successfully")
}
