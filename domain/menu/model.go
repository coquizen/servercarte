package menu

import (
	"errors"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/domain"
)

type SectionType int

const (
	Meal SectionType = iota
	Category
	Container
)

// Section struct defines the service structure.
type Section struct {
	domain.Base
	Title       string      `json:"title" gorm:"unique,not null"`
	Description *string     `json:"description,omitempty"`
	Active      bool        `json:"active" gorm:"default:true"`
	Type        SectionType `json:"type" gorm:"not null, default: 0"`
	Visible     bool        `json:"visible" gorm:"default:true"`
	ListOrder   uint        `json:"list_order" gorm:"default:0"`
	SectionID   *uuid.UUID  `json:"section_id"`
	SubSections []Section   `json:"subsections" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Items       []Item      `json:"items" gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (s *Section) Validate() error {
	if s.Title == "" {
		return errors.New("title is empty")
	}
	if len(s.SubSections) > 0 {
		if len(s.Items) > 0 {
			return errors.New("can only have one of the following: subsections, addons, or items; not both")
		}
	} else if len(s.Items) > 0 {
		if len(s.SubSections) > 0 {
			return errors.New("can only have one of the following: subsections, addons, or items; not both")
		}
	}
	return nil
}

type ItemType int

const (
	Plate ItemType = iota
	Snack
	Side
	AddOn
	Condiment
)

// Item struct defines service items.
type Item struct {
	domain.Base
	Title        string     `json:"title" gorm:"not null"`
	Description  *string    `json:"description"`
	Price        uint64     `json:"price"`
	Active       bool       `json:"active" gorm:"default:true"`
	Type         ItemType   `json:"type" gorm:"default:0"`
	ListOrder    uint       `json:"list_order" gorm:"default:0"`
	SectionID    *uuid.UUID `json:"section_id"`
	AddOnID      *uuid.UUID `json:"add_on_id"`
	AddOn        *Section   `json:"add_on" gorm:"foreignKey:AddOnID"`
	CondimentsID *uuid.UUID `json:"condiments_id"`
	Condiments   *Section   `json:"condiments" gorm:"foreignKey:CondimentsID"`
}

func (i *Item) Validate() error {
	if i.Title == "" {
		return errors.New("item is empty")
	}
	if *i.SectionID == uuid.Nil {
		return errors.New("item must have a section parent")
	}
	return nil
}
