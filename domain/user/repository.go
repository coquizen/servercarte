package user

import (
	"context"
)

// Repository represents the user's interface to its data store
type Repository interface {
	Create(context.Context, *User) error
	List(context.Context) ([]User, error)
	View(context.Context, *User) error
	Search(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, *User) error
}
