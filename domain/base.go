package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model struct to change from using auto-incrementing ID to UUID
type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// BeforeCreate generated a UUID before the creation of the row.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID, err = uuid.NewRandom()
	return err
}