package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"api-shiners/pkg/config"
	"api-shiners/pkg/entities"

	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers() ([]entities.User, error)
	GetUserByID(id uuid.UUID) (entities.User, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// ===== GET ALL USERS (dengan caching Redis opsional) =====
func (s *userService) GetAllUsers() ([]entities.User, error) {
	ctx := context.Background()
	cacheKey := "users:all"

	var users []entities.User

	// 🔹 Coba ambil dari Redis (jika aktif)
	if config.RedisClient != nil {
		val, err := config.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			if err := json.Unmarshal([]byte(val), &users); err == nil {
				return users, nil
			}
		} else if err != nil && err.Error() != "redis: client is closed" {
			log.Println("⚠️ Redis not available or not running, skip caching...")
		}
	}

	// 🔹 Ambil dari DB
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 🔹 Simpan ke cache (jika Redis aktif)
	if config.RedisClient != nil {
		data, _ := json.Marshal(users)
		err := config.RedisClient.Set(ctx, cacheKey, data, 5*time.Minute).Err()
		if err != nil {
			log.Println("⚠️ Failed to cache users:", err)
		}
	}

	return users, nil
}

// ===== GET USER BY ID (dengan caching Redis opsional) =====
func (s *userService) GetUserByID(id uuid.UUID) (entities.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", id.String())

	var user entities.User

	// 🔹 Coba ambil dari Redis (jika aktif)
	if config.RedisClient != nil {
		val, err := config.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			if err := json.Unmarshal([]byte(val), &user); err == nil {
				return user, nil
			}
		} else if err != nil && err.Error() != "redis: client is closed" {
			log.Println("⚠️ Redis not available or not running, skip caching...")
		}
	}

	// 🔹 Ambil dari DB
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return entities.User{}, err
	}

	// 🔹 Simpan ke Redis (jika aktif)
	if config.RedisClient != nil {
		data, _ := json.Marshal(user)
		err := config.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()
		if err != nil {
			log.Println("⚠️ Failed to cache user:", err)
		}
	}

	return user, nil
}
