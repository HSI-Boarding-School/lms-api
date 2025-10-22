package entities

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuizID    uuid.UUID `gorm:"type:uuid;not null"`
	Type      string    `gorm:"type:varchar(50);not null"`
	Text      string    `gorm:"type:text;not null"`
	Position  *int
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`

	// --- Relations ---
	Quiz Quiz `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
}
