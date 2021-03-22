package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/model"
)


// Service represents the user application interface
type Service interface {
	View(context.Context, uuid.UUID) (*model.User, error)
	Find(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	Delete(context.Context, uuid.UUID) error
}
