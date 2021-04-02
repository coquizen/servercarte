package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/domain"
	"github.com/CaninoDev/gastro/server/domain/account"
	"github.com/CaninoDev/gastro/server/domain/user"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type model struct {
	domain.Base
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	User      user.User
	Username  string              `json:"username" gorm:"not null"`
	Password  string              `json:"password" gorm:"not null"`
	Role      account.AccessLevel `json:"role" gorm:"not null"`
	Token     string              `json:"token,omitempty" gorm:"null"`
	LastLogin time.Time           `json:"last_login,omitempty" gorm:"null"`
}
// AccountRepository represents the client to its persistent repository
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository instantiates an instance for data persistence
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (a *AccountRepository) List(ctx context.Context, accounts *[]account.Account) error {
	return a.db.Preload(clause.Associations).Find(&accounts).Error
}
func (a *AccountRepository) Create(ctx context.Context, account *account.Account) error {

	return a.db.Create(&account).Error
}

func (a *AccountRepository) Find(ctx context.Context, account *account.Account) error {
	if account.Username != "" {
		return a.db.First(&account, "username = ?", account.Username).Error
	}
	return a.db.First(&account).Error
}

func (a *AccountRepository) Update(ctx context.Context, account *account.Account) error {
	return a.db.Save(&account).Error
}

func (a *AccountRepository) Delete(ctx context.Context, accountID uuid.UUID) error {
	return a.db.Delete(&account.Account{}, accountID).Error
}



