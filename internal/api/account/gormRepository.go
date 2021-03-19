package account

import (
	"context"

	"gorm.io/gorm"

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

func (a *GormDBRepository) All(ctx context.Context, accounts *[]model.Account) error {
	return a.db.Find(&accounts).Error
}
func (a *GormDBRepository) Create(ctx context.Context, account *model.Account) error {
	return a.db.Create(&account).Error
}

func (a *GormDBRepository) Find(ctx context.Context, account *model.Account) error {
	return a.db.First(&account).Error
}

func (a *GormDBRepository) Update(ctx context.Context, account *model.Account) error {
	return a.db.Save(&account).Error
}

func (a *GormDBRepository) Delete(ctx context.Context, account *model.Account) error {
	return a.db.Delete(&account).Error
}



