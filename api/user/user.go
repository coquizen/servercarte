package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/model"
)

// User is the user's data persistence interface (the above service interface)
type User struct {
	repo Repository
}

func Bind(repo Repository) *User {
	return &User{repo: repo}
}

func Initialize(repo Repository) *User {
	return Bind(repo)
}

// View returns record found by the record's ID
func (u User) View(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	user.ID = id
	if err := u.repo.View(ctx, &user); err != nil {
		return &user, err
	}
	return &user, nil
}

// Find finds a record through its various unique attributes
func (u User) Find(ctx context.Context, user *model.User) error {
	if err := u.repo.Search(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u User) Update(ctx context.Context, user *model.User) error {
	return u.repo.Update(ctx, user)
}

func (u User) Delete(ctx context.Context, id uuid.UUID) error {
	var user model.User
	user.ID = id
	return u.repo.Delete(ctx, &user)
}