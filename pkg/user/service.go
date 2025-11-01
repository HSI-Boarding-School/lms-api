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
	GetAllUsers(page, perPage int) ([]entities.User, int64, error)
	GetUserByID(id uuid.UUID) (entities.User, error)
	SetUserRole(ctx context.Context, userID uuid.UUID, roleName string) (*entities.User, error)
	DeactivateUser(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	ActivateUser(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	GetProfile(ctx context.Context, userID string) (*entities.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) (*entities.User, error)
}

type userService struct {
	userRepo UserRepository
	authRepo auth.AuthRepository
}

func NewUserService(userRepo UserRepository, authRepo auth.AuthRepository) UserService {
	return &userService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}


func (s *userService) GetAllUsers(page, perPage int) ([]entities.User, int64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("users:page:%d:perpage:%d", page, perPage)

	var users []entities.User
	var total int64

	if config.RedisClient != nil {
		val, err := config.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			var cached struct {
				Users []entities.User `json:"users"`
				Total int64           `json:"total"`
			}
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				return cached.Users, cached.Total, nil
			}
		}
	}

	users, total, err := s.userRepo.GetAll(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	if config.RedisClient != nil {
		cached := struct {
			Users []entities.User `json:"users"`
			Total int64           `json:"total"`
		}{Users: users, Total: total}

		data, _ := json.Marshal(cached)
		config.RedisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return users, total, nil
}


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

func (s *userService) SetUserRole(ctx context.Context, userID uuid.UUID, roleName string) (*entities.User, error) {

	if s.userRepo == nil || s.authRepo == nil {
		return nil, errors.New("userRepo atau authRepo belum diinisialisasi dengan benar")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	role, err := s.authRepo.FindRoleByName(ctx, roleName)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if err := s.authRepo.RemoveAllRolesFromUser(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to clear old roles: %v", err)
	}

	userRole := &entities.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}
	if err := s.authRepo.AssignUserRole(ctx, userRole); err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (s *userService) DeactivateUser(ctx context.Context, userID uuid.UUID) (*entities.User, error) {

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.userRepo.DeactivateUser(ctx, userID); err != nil {
		return nil, err
	}

	user.IsActive = false

	return &user, nil
}

func (s *userService) ActivateUser(ctx context.Context, userID uuid.UUID) (*entities.User, error) {

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.userRepo.ActivateUser(ctx, userID); err != nil {
		return nil, err
	}

	user.IsActive = true

	return &user, nil
}


func (s *userService) GetProfile(ctx context.Context, userID string) (*entities.User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(uid)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) (*entities.User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}

	user, err := s.userRepo.UpdateProfile(ctx, userID, name, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}




