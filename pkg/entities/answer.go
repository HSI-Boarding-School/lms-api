package entities

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AttemptID  uuid.UUID `gorm:"type:uuid;not null"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"`
	ChoiceID   *uuid.UUID
	IsCorrect  bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"default:now()"`
	UpdatedAt  time.Time `gorm:"default:now()"`

	Attempt  QuizAttempt `gorm:"foreignKey:AttemptID;constraint:OnDelete:CASCADE"`
	Question Question    `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
	Choice   *Choice     `gorm:"foreignKey:ChoiceID"`
}
