package menu

import (
	"context"
)

// Repository represents the expected methods that a database implementation must have to satisfy the contracts of
// this application
type Repository interface {
	ListSections(context.Context) (*[]Section, error)
	FindSection(context.Context, *Section) error
	CreateSection(context.Context, *Section) error
	UpdateSection(context.Context, *Section) error
	UpdateSectionParent(context.Context, *Section, *Section) error
	DeleteSection(context.Context, *Section) error
	ListItems(context.Context) (*[]Item, error)
	FindItem(context.Context, *Item) error
	CreateItem(context.Context, *Item) error
	UpdateItem(context.Context, *Item) error
	UpdateItemParent(context.Context, *Item, *Section) error
	DeleteItem(context.Context, *Item) error
}


