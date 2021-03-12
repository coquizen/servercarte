package user

import (
	"context"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// Repository represents the user's interface to its data store
type Repository interface {
	Create(context.Context, *model.User) error
 	View(context.Context, *model.User) error
	Search(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	Delete(context.Context, *model.User) error
}
