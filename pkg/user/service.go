package user

import (
	"context"
	"encoding/json"
	"fmt"
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

//
// ===== GET ALL USERS (dengan caching Redis) =====
//
func (s *userService) GetAllUsers() ([]entities.User, error) {
	ctx := context.Background()
	cacheKey := "users:all"

	// 🔹 Cek dari Redis dulu
	val, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && val != "" {
		var users []entities.User
		if err := json.Unmarshal([]byte(val), &users); err == nil {
			return users, nil
		}
	}

	// 🔹 Jika cache kosong → ambil dari DB
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 🔹 Simpan ke cache selama 5 menit
	data, _ := json.Marshal(users)
	config.RedisClient.Set(ctx, cacheKey, data, 5*time.Minute)

	return users, nil
}

//
// ===== GET USER BY ID (dengan caching Redis) =====
//
func (s *userService) GetUserByID(id uuid.UUID) (entities.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", id.String())

	// 🔹 Cek cache Redis
	val, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && val != "" {
		var user entities.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return user, nil
		}
	}

	// 🔹 Jika cache kosong → ambil dari DB
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return entities.User{}, err
	}

	// 🔹 Simpan ke Redis (TTL 10 menit)
	data, _ := json.Marshal(user)
	config.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute)

	return user, nil
}
