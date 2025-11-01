package entities

import (
	"time"

	"github.com/google/uuid"
)

type FeedbackAnswer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null" json:"question_id"`
	StudentID  uuid.UUID `gorm:"type:uuid;not null" json:"student_id"`
	Answer     string    `gorm:"type:text;not null" json:"answer"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`

	Question FeedbackQuestion `gorm:"foreignKey:QuestionID" json:"-"`
	Student  User             `gorm:"foreignKey:StudentID;references:ID" json:"student"`
}
