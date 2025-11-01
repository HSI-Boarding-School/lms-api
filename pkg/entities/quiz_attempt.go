package entities

import (
	"time"

	"github.com/google/uuid"
)

type QuizAttempt struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuizID      uuid.UUID  `gorm:"type:uuid;not null"`
	StudentID   uuid.UUID  `gorm:"type:uuid;not null"`
	StartedAt   time.Time  `gorm:"default:now()"`
	SubmittedAt *time.Time `gorm:"column:submited_at"` // sesuai kolom SQL
	Score       float64    `gorm:"type:numeric(5,2)"`
	DurationSec *int
	CreatedAt   time.Time  `gorm:"default:now()"`
	UpdatedAt   time.Time  `gorm:"default:now()"`

	Quiz    Quiz  `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
	Student User  `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
	Answers []Answer `gorm:"foreignKey:AttemptID"`
}
