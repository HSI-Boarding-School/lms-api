package entities

import (
	"time"

	"github.com/google/uuid"
)

type FeedbackQuestion struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Question  string    `gorm:"type:text;not null" json:"question"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:now()" json:"updated_at"`

	Answers []FeedbackAnswer `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}
