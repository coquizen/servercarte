package menu

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// service describes the expected behavior for manipulating menu data.
type Service interface {
	Sections(context.Context) (*[]model.Section, error)
	SectionByID(context.Context, string) (*model.Section, error)
	NewSection(context.Context, *model.Section) error
	UpdateSectionData(context.Context, *model.Section) error
	ReParentSection(context.Context, *model.Section, uuid.UUID) error
	DeleteSection(context.Context, string) error
	Items(context.Context) (*[]model.Item, error)
	ItemByID(context.Context, string) (*model.Item, error)
	NewItem(context.Context, *model.Item) error
	ReParentItem(context.Context, *model.Item, uuid.UUID) error
	UpdateItemData(context.Context, *model.Item) error
	DeleteItem(context.Context, string) error
}



