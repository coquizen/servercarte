package account

import (
	"context"
	"time"

	"github.com/CaninoDev/gastro/server/api"
	"github.com/CaninoDev/gastro/server/api/user"

	"github.com/google/uuid"
)


type AccessLevel int

const (
	Admin AccessLevel = iota
	Employee
	Guest
)

// Account contains a user's account properties for interacting with the system
type Account struct {
	api.Base
	UserID    uuid.UUID   `json:"userID" gorm:"not null"`
	User      user.User   `comment:"account BELONGS TO user"`
	Username  string      `json:"username" gorm:"not null"`
	Password  string      `json:"password" gorm:"not null"`
	Role      AccessLevel `json:"role" gorm:"not null"`
	Token     string      `json:"token,omitempty" gorm:"null"`
	LastLogin time.Time   `json:"last_login,omitempty" gorm:"null"`
}

// NewAccountRequest represent the request struct for Create endpoint
type NewAccountRequest struct {
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Address1        string  `json:"address_1"`
	Address2        *string `json:"address_2,omitempty"`
	ZipCode         uint    `json:"zip_code"`
	Username        string  `json:"username"`
	Password        string  `json:"password"`
	PasswordConfirm  string  `json:"password_confirm"`
	Email           string  `json:"email"`
}

// NewAccountRequest represent the request struct for Update endpoint
type UpdateAccountRequest struct {
	ID       uuid.UUID `json:"-"`
	Address1 string    `json:"address_1,omitempty"`
	Address2 *string   `json:"address_2,omitempty"`
	ZipCode  uint      `json:"zip_code,omitempty"`
	Email    string    `json:"email,omitempty"`
}

// Service describes the expected behavior in creating accounts and establishing roles for users for the purposes of
// ACL and RBAC
type Service interface {
	New(ctx context.Context, newAccountDetails NewAccountRequest) error
	Authenticate(ctx context.Context, username string, password string) (string, error)
	List(ctx context.Context) ([]Account, error)
	Update(ctx context.Context, request UpdateAccountRequest) error
	Find(ctx context.Context, username string) (*Account, error)
	Delete(ctx context.Context, accountID uuid.UUID, password string) error
}

