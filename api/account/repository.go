package account

import (
	"context"

	"github.com/CaninoDev/gastro/server/api"
)

// Repository describes the expected behavior for the data persistence of
// account information.
type Repository interface {
	All(ctx context.Context, accounts *[]api.Account) error
	Create(ctx context.Context, account *api.Account) error
	Find(ctx context.Context, account *api.Account) error
	Update(ctx context.Context, account *api.Account) error
	Delete(ctx context.Context, account *api.Account) error
}