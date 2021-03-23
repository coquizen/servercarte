package api

import (
	"github.com/google/uuid"
)

// Section struct defines the menu structure.
type Section struct {
	Base
	Title           string     `json:"title" gorm:"unique,not null"`
	Description     *string    `json:"description,omitempty"`
	Active          bool       `json:"active" gorm:"default:true"`
	Visible         bool       `json:"visible" gorm:"default:true"`
	ListOrder       uint      `json:"list_order" gorm:"default:0"`
	SectionID       *uuid.UUID  `json:"section_id,omitempty,"`
	SubSections     []Section `json:"subsections" gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,
OnDelete:SET NULL;"`
	Items           []Item `json:"items" gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,
OnDelete:SET NULL"`
	AddOns          []Item `json:"add_ons,
omitempty" gorm:"foreignKey:AddOnsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
