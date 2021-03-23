package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/api"
)


// Service represents the user application interface
type Service interface {
	View(context.Context, uuid.UUID) (*api.User, error)
	Find(context.Context, *api.User) error
	Update(context.Context, *api.User) error
	Delete(context.Context, uuid.UUID) error
}
