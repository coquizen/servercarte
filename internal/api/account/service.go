package account

import (
	"context"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// Service describes the expected behavior in creating accounts and establishing roles for users for the purposes of
// ACL and RBAC
type Service interface {
	List(ctx context.Context) (*[]model.Account, error)
	Create(ctx context.Context, req *createAccountRequest) error
	Authenticate(ctx context.Context, username string, password string) (string, error)
	Update(ctx context.Context, id string, request *updateAccountRequest) error
	FindByUsername(ctx context.Context, username string) (*model.Account, error)
	FindByToken(ctx context.Context, token string) (*model.Account, error)
	Delete(ctx context.Context, accountID string, password string) error
}


// TODO: implement ACL and Role based access