package account

import (
	"context"
	"errors"
	"time"

	"github.com/CaninoDev/gastro/server/authentication"
	"github.com/CaninoDev/gastro/server/domain/user"
	"github.com/CaninoDev/gastro/server/security"

	"github.com/google/uuid"
)

// Service describes the expected behavior in creating accounts and establishing roles for users for the purposes of
// ACL and RBAC
type Service interface {
	New(ctx context.Context, newAccountDetails NewAccountRequest) error
	Accounts(ctx context.Context) (*[]Account, error)
	Update(ctx context.Context, request UpdateAccountRequest) error
	Find(ctx context.Context, username string) (*Account, error)
	Delete(ctx context.Context, accountID uuid.UUID) error
	Authenticate(ctx context.Context, username, password string) (*Account, error)
}

// service are the contracted methods to interact with GORM
type service struct {
	accountRepo Repository
	userSvc     user.Service
	secSvc      security.Service
	authSvc     authentication.Service
}

// NewService returns a new instance of service
func NewService(accountRepo Repository, userSvc user.Service, secSvc security.Service,
	authSvc authentication.Service) *service {
	return &service{accountRepo, userSvc, secSvc, authSvc}
}

// New creates a new account (bringing in the user model).
// Each user will have zero (restaurant guest) to 1 (
// restaurant employee). This domain is primarily managed by the restaurant's manager
func (a *service) New(ctx context.Context, req NewAccountRequest) error {
	newAccount, newUser := req.unwrap()
	// Ensure that the user doesn't already exist.
	if err := a.userSvc.Find(ctx, newUser); err == nil {
		return ErrUserAlreadyExists
	}

	// Check if password matches and complies with policy
	if !a.secSvc.ConfirmationChecker(ctx, req.Password, req.PasswordConfirm) {
		return errors.New("passwords don't match")
	}
	if err := a.secSvc.IsValid(ctx, req.Password); err != nil {
		return err
	}

	// Make sure username isn't already used by another account
	if _, err := a.Find(ctx, req.Username); err == nil {
		return errors.New("username already in use")
	}
	newAccount.Role = req.Role

	// Hash the given password to secure storage
	newAccount.Password = a.secSvc.Hash(ctx, req.Password)

	// Create the user
	if err := a.userSvc.New(ctx, newUser); err != nil {
		return err
	}
	// Assign the user to the newly created account
	newAccount.User = *newUser

	if err := a.accountRepo.Create(ctx, newAccount); err != nil {
		return err
	}
	return nil
}

func (a *service) Accounts(ctx context.Context) (*[]Account,error) {
	var accounts []Account
	if err := a.accountRepo.List(ctx, &accounts); err != nil {
		return &accounts, err
	}
	return &accounts, nil
}

func (a *service) Find(ctx context.Context, username string) (*Account, error) {
	var acct Account
	acct.Username = username
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return &Account{}, err
	}
	return &acct, nil
}

func (a *service) ChangePassword(ctx context.Context, username, oldPassword, newPassword, confirmNewPassword string) error {
	var acct Account
	acct.Username = username
	if newPassword != confirmNewPassword {
		return errors.New("passwords don't match")
	}

	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return err
	}
	if !a.secSvc.VerifyPasswordMatches(ctx, acct.Password, oldPassword) {
		return errors.New("password incorrect")
	}

	encryptedPW := a.secSvc.Hash(ctx, newPassword)
	acct.Password = encryptedPW
	if err := a.accountRepo.Update(ctx, &acct); err != nil {
		return err
	}
	return nil
}

// Delete will delete the intended account
func (a *service) Delete(ctx context.Context, id uuid.UUID) error {
	return a.accountRepo.Delete(ctx, id)
}

func (a *service) List(ctx context.Context) ([]Account, error) {
	var accounts []Account
	if err := a.accountRepo.List(ctx, &accounts); err != nil {
		return accounts, err
	}
	return accounts, nil
}

func (a *service) Update(ctx context.Context, request UpdateAccountRequest) error {
	var account Account
	account.ID = request.ID
	if request.Role != nil {
		account.Role = *request.Role
	}
	if err := a.accountRepo.Update(ctx, &account); err != nil {
		return err
	}
	updateUser := account.User
	if err := a.userSvc.Find(ctx, &updateUser); err != nil {
		return err
	}

	if err := a.accountRepo.Update(ctx, &account); err != nil {
		return err
	}


	if request.Email != nil {
		updateUser.Email = *request.Email
	}
	if request.Address1 != nil {
		updateUser.Address1 = *request.Address1
	}
	if request.Address2 != nil {
		updateUser.Address2 = *request.Address2
	}
	if request.ZipCode != nil {
		updateUser.ZipCode = *request.ZipCode
	}

	if err := a.userSvc.Update(ctx, &updateUser); err != nil {
		return err
	}
	return nil
}

func (a *service) Authenticate(ctx context.Context, username, password string) (*Account, error) {
	acct, err := a.Find(ctx, username)
	if err != nil {
		return acct, ErrAccountNotFound
	}
	if !a.secSvc.VerifyPasswordMatches(ctx, acct.Password, password) {
		return acct, ErrUnauthorized
	}
	acct.LastLogin = time.Now().UTC()
	if err := a.accountRepo.Update(ctx, acct); err != nil {
		return acct, err
	}
	return acct, nil
}
