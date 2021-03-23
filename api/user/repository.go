package user

import (
	"context"

	"github.com/CaninoDev/gastro/server/api"
)

// Repository represents the user's interface to its data store
type Repository interface {
	Create(context.Context, *api.User) error
 	View(context.Context, *api.User) error
	Search(context.Context, *api.User) error
	Update(context.Context, *api.User) error
	Delete(context.Context, *api.User) error
}
