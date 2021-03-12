package menu

import (
	"context"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// Repository represents the expected methods that a database implementation must have to satisfy the contracts of
// this application
type Repository interface {
	ListSections(context.Context) (*[]model.Section, error)
	FindSection(context.Context, *model.Section) error
	CreateSection(context.Context, *model.Section) error
	UpdateSection(context.Context, *model.Section) error
	UpdateSectionParent(context.Context, *model.Section, *model.Section) error
	DeleteSection(context.Context, *model.Section) error
	ListItems(context.Context) (*[]model.Item, error)
	FindItem(context.Context, *model.Item) error
	CreateItem(context.Context, *model.Item) error
	UpdateItem(context.Context, *model.Item) error
	UpdateItemParent(context.Context, *model.Item, *model.Section) error
	DeleteItem(context.Context, *model.Item) error
}
