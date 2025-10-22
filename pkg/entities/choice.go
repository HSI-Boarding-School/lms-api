package entities

import (
	"time"

	"github.com/google/uuid"
)

type Choice struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"`
	Text       string    `gorm:"type:text;not null"`
	IsCorrect  bool      `gorm:"default:false"`
	Position   *int
	CreatedAt  time.Time `gorm:"default:now()"`
	UpdatedAt  time.Time `gorm:"default:now()"`

	// --- Relations ---
	Question Question `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}
