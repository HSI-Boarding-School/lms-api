package entities

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code           string    `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Title          string    `gorm:"size:255;not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description,omitempty"`
	OwnerTeacherID uuid.UUID `gorm:"type:uuid" json:"owner_teacher_id"`
	IsPublished    bool      `gorm:"default:false" json:"is_published"`
	CreatedAt      time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:now()" json:"updated_at"`

	OwnerTeacher *User `gorm:"foreignKey:OwnerTeacherID;constraint:OnDelete:SET NULL;" json:"owner_teacher,omitempty"`

	Modules []CourseModule `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE;" json:"modules,omitempty"`
}
