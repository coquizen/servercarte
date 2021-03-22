package account

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// NewAccountRequest represent the request struct for Create endpoint
type NewAccountRequest struct {
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Address1        string  `json:"address_1"`
	Address2        *string `json:"address_2,omitempty"`
	ZipCode         uint    `json:"zip_code"`
	Username        string  `json:"username"`
	Password        string  `json:"password"`
	PasswordConfirm string  `json:"password_confirm"`
	Email           string  `json:"email"`
}

// NewAccountRequest represent the request struct for Update endpoint
type UpdateAccountRequest struct {
	ID          uuid.UUID `json:"-"`
	Address1        string  `json:"address_1,omitempty"`
	Address2        *string `json:"address_2,omitempty"`
	ZipCode         uint    `json:"zip_code,omitempty"`
	Email           string  `json:"email,omitempty"`
}

// Service describes the expected behavior in creating accounts and establishing roles for users for the purposes of
// ACL and RBAC
type Service interface {
	New(ctx context.Context, newAccountDetails NewAccountRequest) error
	Authenticate(ctx context.Context, username string, password string) (string, error)
	List(ctx context.Context) (*[]model.Account, error)
	Update(ctx context.Context, accountID uuid.UUID, request UpdateAccountRequest) error
	FindByUsername(ctx context.Context, username string) (*model.Account, error)
	FindByToken(ctx context.Context, token string) (*model.Account, error)
	Delete(ctx context.Context, accountID uuid.UUID, password string) error
}

// TODO: implement ACL and Role based access
