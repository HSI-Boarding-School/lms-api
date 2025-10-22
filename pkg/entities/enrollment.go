package entities

import (
	"time"

	"github.com/google/uuid"
)

// Enum di sisi Go
type CourseRole string

const (
	CourseRoleStudent CourseRole = "STUDENT"
	CourseRoleTeacher CourseRole = "TEACHER"
)

// Entity Enrollment
type Enrollment struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null"`
	CourseID     uuid.UUID  `gorm:"type:uuid;not null"`
	RoleInCourse CourseRole `gorm:"type:varchar(50);not null"` // pakai VARCHAR
	EnrolledAt   time.Time  `gorm:"default:now()"`

	// --- Relations ---
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}
