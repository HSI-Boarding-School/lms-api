package entities

import (
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ModuleID        uuid.UUID  `gorm:"type:uuid;not null"`
	Title           string     `gorm:"type:varchar(255);not null"`
	Instructions    string     `gorm:"type:text"`
	OpenAt          *time.Time `gorm:"type:timestamp"`
	CloseAt         *time.Time `gorm:"type:timestamp"`
	TimeLimitSec    *int       `gorm:"column:time_limit_sec"`
	AttemptAllowed  int        `gorm:"default:1"`
	IsPublished     bool       `gorm:"default:false"`
	CreatedAt       time.Time  `gorm:"default:now()"`
	UpdatedAt       time.Time  `gorm:"default:now()"`

	Module    CourseModule  `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE"`
	Questions []Question    `gorm:"foreignKey:QuizID"`
}
