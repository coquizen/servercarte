package gormDB

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/CaninoDev/gastro/server/api"
)

// AccountRepository represents the client to its persistent storage
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository instantiates an instance for data persistence
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (a *AccountRepository) All(ctx context.Context, accounts *[]api.Account) error {
	return a.db.Preload(clause.Associations).Find(&accounts).Error
}
func (a *AccountRepository) Create(ctx context.Context, account *api.Account) error {
	return a.db.Create(&account).Error
}

func (a *AccountRepository) Find(ctx context.Context, account *api.Account) error {
	if account.Username != "" {
		return a.db.First(&account, "username = ?", account.Username).Error
	}
	return a.db.First(&account).Error
}

func (a *AccountRepository) Update(ctx context.Context, account *api.Account) error {
	return a.db.Save(&account).Error
}

func (a *AccountRepository) Delete(ctx context.Context, account *api.Account) error {
	return a.db.Delete(&account).Error
}



