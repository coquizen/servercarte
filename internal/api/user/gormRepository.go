package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/CaninoDev/gastro/server/internal/logger"
	"github.com/CaninoDev/gastro/server/internal/model"
)

// GormDBRepository represents the client to its persistent storage
type GormDBRepository struct {
	db *gorm.DB
}

// NewGormDBRepository instantiates an instance for data persistence
func NewGormDBRepository(db *gorm.DB) *GormDBRepository {
	return &GormDBRepository{db}
}

func (r GormDBRepository) View(ctx context.Context, user *model.User) error {
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

func (r GormDBRepository) Search(ctx context.Context, user *model.User) error {
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

func (r GormDBRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.Create(&user).Error
}
func (r GormDBRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.Save(&user).Error
}

func (r GormDBRepository) Delete(ctx context.Context, user *model.User) error {

	return r.db.Delete(&model.User{}, "id = ?", user.ID).Error
}
