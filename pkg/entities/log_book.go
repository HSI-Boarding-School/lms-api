package entities

import (
	"time"

	"github.com/google/uuid"
)

type StatusType string

const (
	StatusDraft     StatusType = "DRAFT"
	StatusSubmitted StatusType = "SUBMITTED"
	StatusLocked    StatusType = "LOCKED"
)

type LogBook struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseID     uuid.UUID  `gorm:"type:uuid;not null"`
	StudentID    uuid.UUID  `gorm:"type:uuid;not null;column:student_id_id"`
	PeriodStart  time.Time  `gorm:"type:date;not null"`
	PeriodEnd    time.Time  `gorm:"type:date;not null"`
	Status       StatusType `gorm:"type:varchar(20);default:'DRAFT'"`
	SubmittedAt  *time.Time `gorm:"column:submited_at"`
	LockedAt     *time.Time
	CreatedAt    time.Time  `gorm:"default:now()"`
	UpdatedAt    time.Time  `gorm:"default:now()"`

	Course  Course  `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	Student User    `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
	Entries []LogBookEntry `gorm:"foreignKey:LogBookID"`
}
