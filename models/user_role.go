package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RoleID     uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	AssignedAt time.Time `gorm:"default:now()" json:"assigned_at"`
	AssignedBy uuid.UUID `gorm:"type:uuid" json:"assigned_by"`
}
