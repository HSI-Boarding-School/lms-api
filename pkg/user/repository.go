package user

import (
	"api-shiners/pkg/entities"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll(page, perPage int) ([]entities.User, int64, error)
	GetByID(id uuid.UUID) (entities.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) (*entities.User, error)
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
	ActivateUser(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll(page, perPage int) ([]entities.User, int64, error) {
	var users []entities.User
	var total int64

	query := r.db.Model(&entities.User{}).Preload("Roles")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Limit(perPage).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetByID(id uuid.UUID) (entities.User, error) {
	var user entities.User
	err := r.db.Preload("Roles").First(&user, "id = ?", id).Error
	return user, err
}

func (r *userRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	user.Name = name
	user.Email = email

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}


func (r *userRepository) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", userID).
		Update("is_active", false).Error
}

func (r *userRepository) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", userID).
		Update("is_active", true).Error
}






