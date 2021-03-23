package menu

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/api"
)

// Service describes the expected behavior for manipulating menu data.
type Service interface {
	Sections(context.Context) (*[]api.Section, error)
	SectionByID(context.Context, string) (*api.Section, error)
	NewSection(context.Context, *api.Section) error
	UpdateSectionData(context.Context, *api.Section) error
	ReParentSection(context.Context, *api.Section, uuid.UUID) error
	DeleteSection(context.Context, string) error
	Items(context.Context) (*[]api.Item, error)
	ItemByID(context.Context, string) (*api.Item, error)
	NewItem(context.Context, *api.Item) error
	ReParentItem(context.Context, *api.Item, uuid.UUID) error
	UpdateItemData(context.Context, *api.Item) error
	DeleteItem(context.Context, string) error
}



