package menu

import (
	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/api"
)

// Section struct defines the service structure.
type Section struct {
	api.Base
	Title           string     `json:"title" gorm:"unique,not null"`
	Description     *string    `json:"description,omitempty"`
	Active          bool       `json:"active" gorm:"default:true"`
	Visible         bool       `json:"visible" gorm:"default:true"`
	ListOrder       uint      `json:"list_order" gorm:"default:0"`
	SectionID       *uuid.UUID  `json:"section_id,omitempty,"`
	SubSections     []Section   `json:"subsections" gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Items           []Item `json:"items" gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AddOns          []Item `json:"add_ons,omitempty" gorm:"foreignKey:AddOnsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Item struct defines service items.
type Item struct {
	api.Base
	Title       string    `json:"title" gorm:"not null"`
	Description *string   `json:"description"`
	Price       uint64    `json:"price"`
	Active      bool      `json:"active" gorm:"default:true"`
	Type        uint      `json:"type,omitempty" gorm:"default:1"`
	ListOrder   uint      `json:"list_order" gorm:"default:0"`
	SectionID   *uuid.UUID `json:"section_id"`
	AddOnsID    *uuid.UUID

}

