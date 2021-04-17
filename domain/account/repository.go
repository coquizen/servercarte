package account

import (
	"context"

	"github.com/google/uuid"
)

// Repository describes the expected behavior for the data persistence of
// account information.
type Repository interface {
	List(ctx context.Context) ([]Account, error)
	Create(ctx context.Context, account *Account) error
	Find(ctx context.Context, account *Account) error
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, accountID uuid.UUID) error
}
