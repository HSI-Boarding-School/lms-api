package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaterialType string

const (
	MaterialTypeFile MaterialType = "FILE"
	MaterialTypeLink MaterialType = "LINK"
	MaterialTypeText MaterialType = "TEXT"
)

type Material struct {
	ID           uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ModuleID     uuid.UUID    `gorm:"type:uuid;not null" json:"module_id"`
	Module       CourseModule `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"module"`
	Type         MaterialType `gorm:"type:varchar(50);not null" json:"type"`
	Title        string       `gorm:"type:varchar(255);not null" json:"title"`
	ContentOrURL string       `gorm:"type:text" json:"content_or_url"`
	FileKey      string       `gorm:"type:varchar(255)" json:"file_key"`
	Position     int          `json:"position"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
