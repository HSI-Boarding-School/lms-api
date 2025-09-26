// services/auth_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/daffa-fawwaz/shiners-lms-backend/app/models"
	"github.com/daffa-fawwaz/shiners-lms-backend/app/repositories"
	"github.com/daffa-fawwaz/shiners-lms-backend/app/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, string, string, error) // access, refresh, role
}

type authService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// Register user baru dengan role default = STUDENT
func (s *authService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	// cek apakah email sudah ada
	existing, _ := s.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// buat user baru
	user := &models.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// assign role default
	role, _ := s.roleRepo.FindByName(ctx, "STUDENT")
	if role != nil {
		_ = s.userRepo.AssignRole(ctx, user.ID, role.ID, nil)
	}

	return user, nil
}

// Login user, return access & refresh token + role
func (s *authService) Login(ctx context.Context, email, password string) (string, string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", "", errors.New("invalid email or password")
	}

	// cek password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", "", errors.New("invalid email or password")
	}

	// ambil role user
	role, err := s.roleRepo.FindByUserID(ctx, user.ID.String())
	if err != nil || role == nil {
		return "", "", "", errors.New("failed to fetch user role")
	}

	// generate tokens
	accessToken, err := utils.GenerateJWT(user.ID.String(), 15*time.Minute)
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), 7*24*time.Hour)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, role.Name, nil
}
