package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoleName string

const (
	ADMIN RoleName = "ADMIN"
	TEACHER RoleName = "TEACHER"
	STUDENT RoleName = "STUDENT"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name 		RoleName `gorm:"type:varchar(20);default:'STUDENT'" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`

	Users []*User `gorm:"many2many:user_roles;" json:"-"`
}
