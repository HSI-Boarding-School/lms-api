package utils

import (
	"github.com/daffa-fawwaz/shiners-lms-backend/dto"
	"github.com/daffa-fawwaz/shiners-lms-backend/models"
)

func ToUserResponse(user *models.User) dto.UserResponse {
	roles := make([]dto.RoleResponse, 0)
	for _, r := range user.Roles {
		roles = append(roles, dto.RoleResponse{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
		Roles:    roles,
	}
}
