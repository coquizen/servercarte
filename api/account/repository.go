package account

import (
	"context"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// Repository describes the expected behavior for the data persistence of
// account information.
type Repository interface {
	All(ctx context.Context, accounts *[]model.Account) error
	Create(ctx context.Context, account *model.Account) error
	Find(ctx context.Context, account *model.Account) error
	Update(ctx context.Context, account *model.Account) error
	Delete(ctx context.Context, account *model.Account) error
}