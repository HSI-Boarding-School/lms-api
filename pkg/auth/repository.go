package auth

import (
	"api-shiners/pkg/entities"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	FindRoleByName(ctx context.Context, name string) (*entities.Role, error)
	AssignUserRole(ctx context.Context, userRole *entities.UserRole) error
	UpdateUser(ctx context.Context, user *entities.User) error

	// 🔹 Tambahan untuk Forgot & Reset Password
	SaveResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	FindByResetToken(ctx context.Context, token string) (*entities.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error
	ClearResetToken(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// ================================================================
// BASIC USER CRUD
// ================================================================

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Preload("Roles").
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindRoleByName(ctx context.Context, name string) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) AssignUserRole(ctx context.Context, userRole *entities.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// ================================================================
// FORGOT & RESET PASSWORD (REAL DYNAMIC VERSION)
// ================================================================

func (r *userRepository) SaveResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"reset_token":   token,
			"reset_expires": expiresAt,
		}).Error
}

func (r *userRepository) FindByResetToken(ctx context.Context, token string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).
		Where("reset_token = ?", token).
		First(&user).Error; err != nil {
		return nil, err
	}

	if user.ResetExpires == nil {
		return nil, errors.New("invalid reset token")
	}

	if time.Now().After(*user.ResetExpires) {
		return nil, errors.New("reset token expired")
	}

	return &user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error {
	return r.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Update("password_hash", newPasswordHash).Error
}

func (r *userRepository) ClearResetToken(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"reset_token":   nil,
			"reset_expires": nil,
		}).Error
}
