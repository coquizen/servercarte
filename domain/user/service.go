package user

import (
	"context"

	"github.com/google/uuid"
)

// Service represents the user application interface
type Service interface {
	Create(context.Context, *User) error
	Users(context.Context) ([]User, error)
	View(context.Context, uuid.UUID) (*User, error)
	Find(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, uuid.UUID) error
}

// service is the user's data persistence interface (the Service interface)
type service struct {
	repo Repository
}

// NewService instantiates an user service and returns an instance containing its persistent store
func NewService(userRepo Repository) *service {
	return &service{
		userRepo,
	}
}

// New creates a new user. This user can be either a restaurant guest, employee, or manager/owner.
func (u *service) Create(ctx context.Context, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return u.repo.Create(ctx, user)
}

// Users returns all the users currently in the system
func (u *service) Users(ctx context.Context) ([]User, error) {
	return u.repo.List(ctx)
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
	if err := user.Validate(); err != nil {
		return err
	}
	return u.repo.Update(ctx, user)
}

func (u *service) Delete(ctx context.Context, id uuid.UUID) error {
	var user User
	user.ID = id
	return u.repo.Delete(ctx, &user)
}
