package menu

import (
	"context"

	"github.com/CaninoDev/gastro/server/api"
)

// Repository represents the expected methods that a database implementation must have to satisfy the contracts of
// this application
type Repository interface {
	ListSections(context.Context) (*[]api.Section, error)
	FindSection(context.Context, *api.Section) error
	CreateSection(context.Context, *api.Section) error
	UpdateSection(context.Context, *api.Section) error
	UpdateSectionParent(context.Context, *api.Section, *api.Section) error
	DeleteSection(context.Context, *api.Section) error
	ListItems(context.Context) (*[]api.Item, error)
	FindItem(context.Context, *api.Item) error
	CreateItem(context.Context, *api.Item) error
	UpdateItem(context.Context, *api.Item) error
	UpdateItemParent(context.Context, *api.Item, *api.Section) error
	DeleteItem(context.Context, *api.Item) error
}
