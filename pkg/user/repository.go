package user

import (
	"api-shiners/pkg/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll() ([]entities.User, error)
	GetByID(id uuid.UUID) (entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll() ([]entities.User, error) {
	var users []entities.User
	err := r.db.Preload("Roles").Find(&users).Error
	return users, err
}

func (r *userRepository) GetByID(id uuid.UUID) (entities.User, error) {
	var user entities.User
	err := r.db.Preload("Roles").First(&user, "id = ?", id).Error
	return user, err
}




