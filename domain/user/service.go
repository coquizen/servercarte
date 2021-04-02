package user

import (
	"context"

	"github.com/google/uuid"
)


// UserServicee represents the user application interface
type Service interface {
	New(context.Context, *User) error
	View(context.Context, uuid.UUID) (*User, error)
	Find(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, uuid.UUID) error
}

// service is the user's data persistence interface (the above accountservice interface)
type service struct {
	repo Repository
}

// NewService constructs a new accountservice accountservice and returns an instance containing its persistent store
func NewService(userRepo Repository) *service {
	return &service{
		userRepo,
	}
}
//
// type NewUserRequest struct {
// 	Name            string
// 	Address1        string
// 	Address2        *string
// 	ZipCode         int
// 	TelephoneNumber string
// 	Email           string
// }
// New creates a new user. This user can be either a restaurant guest, employee, or manager/owner.
func (u *service) New(ctx context.Context, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return u.repo.Create(ctx, user)
}

// View returns record found by the record's ID
func (u *service) View(ctx context.Context, id uuid.UUID) (*User, error) {
	var currentUser User
	currentUser.ID = id
	if err := u.repo.View(ctx, &currentUser); err != nil {
		return &currentUser, err
	}
	return &currentUser, nil
}

// Find finds a record through its various unique attributes
func (u *service) Find(ctx context.Context, user *User) error {
	if err := u.repo.Search(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *service) Update(ctx context.Context, user *User) error {
	return u.repo.Update(ctx, user)
}

func (u *service) Delete(ctx context.Context, id uuid.UUID) error {
	var user User
	user.ID = id
	return u.repo.Delete(ctx, &user)
}