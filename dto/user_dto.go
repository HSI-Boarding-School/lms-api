package dto

import (
	"github.com/google/uuid"
)

type RoleResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserResponse struct {
	ID       uuid.UUID     `json:"id"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	IsActive bool          `json:"is_active"`
	Roles    []RoleResponse `json:"roles,omitempty"`
}
