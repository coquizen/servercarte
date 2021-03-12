package model

import (
	"github.com/google/uuid"
)

// Item struct defines menu items.
type Item struct {
	Base
	Title       string    `json:"title" gorm:"not null"`
	Description *string   `json:"description"`
	Price       uint64    `json:"price"`
	Active      bool      `json:"active" gorm:"default:true"`
	Type        uint      `json:"type,omitempty" gorm:"default:1"`
	ListOrder   uint      `json:"list_order" gorm:"default:0"`
	SectionID   *uuid.UUID `json:"section_id"`
	AddOnsID    *uuid.UUID

}
