package gormDB

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/CaninoDev/gastro/server/internal/logger"
	"github.com/CaninoDev/gastro/server/internal/model"
)

// MenuRepository represents the client to its persistent storage
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository instantiates an instance for data persistence
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r UserRepository) View(ctx context.Context, user *model.User) error {
	if err := r.db.First(&user, user.ID).Error; errors.Is(
		err,
		gorm.ErrRecordNotFound) {
		return err
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return nil
}

func (r UserRepository) Search(ctx context.Context, user *model.User) error {
	if err := r.db.Where("email = ? ", user.Email).Or("first_name = ? AND last_name = ?",
		user.FirstName, user.LastName).Or("id = ?", user.ID).First(&user).Error; errors.Is(
		err,
		gorm.ErrRecordNotFound) {
		return err
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return nil
}

func (r UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.Create(&user).Error
}
func (r UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.Save(&user).Error
}

func (r UserRepository) Delete(ctx context.Context, user *model.User) error {

	return r.db.Delete(&model.User{}, "id = ?", user.ID).Error
}
