package menu

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/coquizen/servercarte/domain"
)

//go:generate stringer -type=SectionType
type SectionType int

const (
	UndefinedSection SectionType = iota
	Meal
	Category
	Container
)

func (s SectionType) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *SectionType) UnmarshalText(text []byte) error {
	*s = SectionTypeFromText(string(text))
	return nil
}

// Section struct defines the service structure.
type Section struct {
	domain.Base
	Title       string      `json:"title" gorm:"unique,not null"`
	Description *string     `json:"description"`
	Active      bool        `json:"active" gorm:"default:true"`
	Type        SectionType `json:"type" gorm:"not null, default: 0"`
	Visible     bool        `json:"visible" gorm:"default:true"`
	ListOrder   uint        `json:"list_order" gorm:"default:0"`
	SectionID   *uuid.UUID  `json:"section_id"`
	SubSections []Section   `json:"subsections" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Items       []Item      `json:"items" gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AddOnsID     *uuid.UUID `json:"add_ons_id"`
	CondimentsID *uuid.UUID `json:"condiments_id"`
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

//go:generate stringer -type=ItemType
type ItemType int

const (
	UndefinedItem ItemType = iota
	Plate
	Snack
	Side
	AddOn

	Condiment
)

func (i ItemType) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

func (i *ItemType) UnmarshalText(text []byte) error {
	*i = ItemTypeFromText(string(text))
	return nil
}

// Item struct defines service items.
type Item struct {
	domain.Base
	Title        string     `json:"title" gorm:"not null"`
	Description  *string    `json:"description"`
	Price        uint64     `json:"price" gorm:"default:000"`
	Active       bool       `json:"active" gorm:"default:true"`
	Type         ItemType   `json:"type" gorm:"default:0"`
	ListOrder    uint       `json:"list_order" gorm:"default:0"`
	SectionID    *uuid.UUID `json:"section_id"`
	AddOns       Section    `json:"add_ons" gorm:"foreignKey:AddOnsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Condiments   Section    `json:"condiments" gorm:"foreignKey:CondimentsID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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

func SectionTypeFromText(text string) SectionType {
	switch strings.ToLower(text) {
	case "meal":
		return Meal
	case "category":
		return Category
	case "container":
		return Container
	default:
		return UndefinedSection
	}
}
func ItemTypeFromText(text string) ItemType {
	switch strings.ToLower(text) {
	case "plate":
		return Plate
	case "snack":
		return Snack
	case "side":
		return Side
	case "add_on":
		return AddOn
	case "condiment":
		return Condiment
	default:
		return UndefinedItem
	}

}
