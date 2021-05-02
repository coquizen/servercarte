package gorm

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/domain/account"
	"github.com/CaninoDev/gastro/server/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AccountRepository represents the client to its persistent repository
type AccountRepository struct {
	db *gorm.DB
}

var nullAccount = account.Account{}

// NewAccountRepository instantiates an instance for data persistence
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (a *AccountRepository) List(ctx context.Context) ([]account.Account, error) {
	var accounts []account.Account
	if err := a.db.Preload(clause.Associations).Find(&accounts).Error; err != nil {
		return []account.Account{}, err
	}
	return accounts, nil
}

// Create creates an account with an associated and newly created user
func (a *AccountRepository) Create(ctx context.Context, account *account.Account, user *user.User) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&account).Error; err != nil {
			// return any error will rollback
			return err
		}

		if err := tx.Model(&account).Association("User").Replace(user); err != nil {
			return err
		}
		return nil
	})
}

func (a *AccountRepository) Find(ctx context.Context, username string) (account.Account, error) {
	var account account.Account
	if err := a.db.Where("username = ?", username).First(&account).Error; err != nil {
		return nullAccount, err
	}
	return account, nil
}

func (a *AccountRepository) Update(ctx context.Context, account *account.Account) error {
	return a.db.Save(&account).Error
}

func (a *AccountRepository) Delete(ctx context.Context, accountID uuid.UUID) error {
	return a.db.Delete(&account.Account{}, accountID).Error
}
