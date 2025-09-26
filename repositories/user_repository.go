package repositories

import (
	"context"
	"time"

	"github.com/daffa-fawwaz/shiners-lms-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	AssignRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, assignedBy *uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) AssignRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, assignedBy *uuid.UUID) error {
	userRole := models.UserRole{
		ID:        uuid.New(),
		UserID:    userID,
		RoleID:    roleID,
		AssignedAt: time.Now(),
	}
	if assignedBy != nil {
		userRole.AssignedBy = *assignedBy
	}
	return r.db.WithContext(ctx).Create(&userRole).Error
}
