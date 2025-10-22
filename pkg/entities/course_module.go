package entities

import (
	"time"

	"github.com/google/uuid"
)

type CourseModule struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;index" json:"course_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	OrderNo   int       `json:"order_no"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:now()" json:"updated_at"`

	// Relasi ke Course (belongs to)
	Course *Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE;" json:"-"`
}
