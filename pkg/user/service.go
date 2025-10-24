package user

import (
	"api-shiners/pkg/auth"
	"api-shiners/pkg/config"
	"api-shiners/pkg/entities"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers() ([]entities.User, error)
	GetUserByID(id uuid.UUID) (entities.User, error)
	SetUserRole(ctx context.Context, userID uuid.UUID, roleName string) (*entities.User, error)
}

type userService struct {
	userRepo UserRepository
	authRepo auth.AuthRepository
}

// ✅ Constructor tunggal — wajib dipakai
func NewUserService(userRepo UserRepository, authRepo auth.AuthRepository) UserService {
	return &userService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

// ===== GET ALL USERS (dengan caching Redis opsional) =====
func (s *userService) GetAllUsers() ([]entities.User, error) {
	ctx := context.Background()
	cacheKey := "users:all"

	var users []entities.User

	// 🔹 Coba ambil dari Redis
	if config.RedisClient != nil {
		val, err := config.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			if err := json.Unmarshal([]byte(val), &users); err == nil {
				return users, nil
			}
		}
	}

	// 🔹 Ambil dari DB
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 🔹 Simpan ke Redis
	if config.RedisClient != nil {
		data, _ := json.Marshal(users)
		config.RedisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return users, nil
}

// ===== GET USER BY ID =====
func (s *userService) GetUserByID(id uuid.UUID) (entities.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", id.String())

	var user entities.User

	if config.RedisClient != nil {
		val, err := config.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			if err := json.Unmarshal([]byte(val), &user); err == nil {
				return user, nil
			}
		}
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return entities.User{}, err
	}

	if config.RedisClient != nil {
		data, _ := json.Marshal(user)
		config.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return user, nil
}

// ===== SET USER ROLE =====
func (s *userService) SetUserRole(ctx context.Context, userID uuid.UUID, roleName string) (*entities.User, error) {
	// Cek dependency dulu
	if s.userRepo == nil || s.authRepo == nil {
		return nil, errors.New("userRepo atau authRepo belum diinisialisasi dengan benar")
	}

	// 1. Cari user berdasarkan ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Cari role berdasarkan nama
	role, err := s.authRepo.FindRoleByName(ctx, roleName)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// 3. Hapus semua role lama user
	if err := s.authRepo.RemoveAllRolesFromUser(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to clear old roles: %v", err)
	}

	// 4. Tambahkan role baru
	userRole := &entities.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}
	if err := s.authRepo.AssignUserRole(ctx, userRole); err != nil {
		return nil, err
	}

	// 5. Ambil ulang user untuk dikembalikan (dengan role terbaru)
	updatedUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

