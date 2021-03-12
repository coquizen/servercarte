package account

import (
	"context"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// Repository describes the expected behavior for the data persistence of
// account information.
type Repository interface {
	Create(context.Context, *model.Account) error
	Find(context.Context, *model.Account) error
	Update(context.Context, *model.Account) error
	Delete(context.Context, *model.Account) error
}