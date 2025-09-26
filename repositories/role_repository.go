// repositories/role_repository.go
package repositories

import (
	"context"

	"github.com/daffa-fawwaz/shiners-lms-backend/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindByUserID(ctx context.Context, userID string) (*models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ambil role berdasarkan user_id (join ke tabel user_roles)
func (r *roleRepository) FindByUserID(ctx context.Context, userID string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Select("roles.*").
		Joins("join user_roles ur on ur.role_id = roles.id").
		Where("ur.user_id = ?", userID).
		First(&role).Error

	if err != nil {
		return nil, err
	}
	return &role, nil
}
