package entities

import (
	"time"

	"github.com/google/uuid"
)

type LogBookEntry struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	LogBookID  uuid.UUID `gorm:"type:uuid;not null"`
	EntryDate  time.Time `gorm:"type:date;not null"`
	Content    string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"default:now()"`
	UpdatedAt  time.Time `gorm:"default:now()"`

	LogBook LogBook `gorm:"foreignKey:LogBookID;constraint:OnDelete:CASCADE"`
}
