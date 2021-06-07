package account

import (
	"context"
	"errors"
	"time"

	"github.com/CaninoDev/gastro/server/domain/security"
	"github.com/CaninoDev/gastro/server/internal/helpers"

	"github.com/CaninoDev/gastro/server/domain/authentication"
	"github.com/CaninoDev/gastro/server/domain/user"
	"github.com/google/uuid"
)

var (
	NullAccount = Account{}

	ErrInvalidEmail  = errors.New("email is invalid")
	ErrUsernameInUse = errors.New("username is already in use")
)

// Service describes the expected behavior in creating accounts and establishing roles for users for the purposes of
// ACL and RBAC
type Service interface {
	New(ctx context.Context, newAccountDetails NewAccountRequest) (*Account, error)
	Accounts(ctx context.Context) ([]Account, error)
	Update(ctx context.Context, request UpdateAccountRequest) error
	Find(ctx context.Context, username string) (Account, error)
	Delete(ctx context.Context, accountID uuid.UUID) error
	Authenticate(ctx context.Context, username, password string) (Account, error)
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
	authSvc authentication.Service) Service {
	return &service{accountRepo, userSvc, secSvc, authSvc}
}

// New creates a new account (bringing in the user model).
// Each user will have zero (restaurant guest) to 1 (
// restaurant employee). This domain is primarily managed by the restaurant's manager
func (a *service) New(ctx context.Context, req NewAccountRequest) (*Account, error) {
	newAccount, newUser := req.unwrap()
	if err := a.validateAccountRequest(req, &newUser); err != nil {
		return &NullAccount, err
	}

	if err := a.checkPreexisting(ctx, &newUser, newAccount.Username); err != nil {
		return &NullAccount, ErrUsernameInUse
	}

	if err := a.secSvc.IsValid(req.Password); err != nil {
		return &NullAccount, err
	}

	if err := a.secSvc.ConfirmationChecker(req.Password, req.PasswordConfirm); err != nil {
		return &NullAccount, err
	}

	newAccount.Password = a.secSvc.Hash(req.Password)

	if err := a.userSvc.Create(ctx, &newUser); err != nil {
		return &NullAccount, err
	}

	if err := a.accountRepo.Create(ctx, &newAccount, &newUser); err != nil {
		return &NullAccount, err
	}
	return &newAccount, nil
}

func (a *service) validateAccountRequest(req NewAccountRequest, newUser *user.User) error {
	// Check to see if the email is in the right format
	if !helpers.IsEmailFormat(newUser.Email) {
		return ErrInvalidEmail
	} else {
		return a.checkPassword(req.Password, req.PasswordConfirm)
	}
}

func (a *service) checkPassword(password, passwordConfirm string) error {
	// Check if password matches and complies with policy
	if err := a.secSvc.ConfirmationChecker(password, passwordConfirm); err != nil {
		return err
	} else {
		return a.secSvc.IsValid(password)
	}
}

func (a *service) checkPreexisting(ctx context.Context, newUser *user.User, username string) error {
	if err := a.userSvc.Find(ctx, newUser); err == nil {
		return user.ErrUserAlreadyExists
	}
	if _, err := a.Find(ctx, username); err == nil {
		return ErrUsernameInUse
	}
	return nil
}

func (a *service) Accounts(ctx context.Context) ([]Account, error) {
	return a.accountRepo.List(ctx)
}

func (a *service) Find(ctx context.Context, username string) (Account, error) {
	account, err := a.accountRepo.Find(ctx, username)
	if err != nil {
		return NullAccount, err
	}
	return account, nil
}

func (a *service) ChangePassword(ctx context.Context, username, oldPassword, newPassword, confirmNewPassword string) error {
	if err := a.secSvc.ConfirmationChecker(newPassword, confirmNewPassword); err != nil {
		return err
	}

	acct, err := a.Authenticate(ctx, username, oldPassword)
	if err != nil {
		return err
	}
	acct.Password = a.secSvc.Hash(newPassword)

	return a.accountRepo.Update(ctx, &acct)
}

// Delete will delete the intended account
func (a *service) Delete(ctx context.Context, id uuid.UUID) error {
	return a.accountRepo.Delete(ctx, id)
}

func (a *service) List(ctx context.Context) ([]Account, error) {
	return a.accountRepo.List(ctx)
}

func (a *service) Update(ctx context.Context, request UpdateAccountRequest) error {
	updatingAccount, updatingUser := request.unwrap()

	if err := a.accountRepo.Update(ctx, updatingAccount); err != nil {
		return err
	}
	updatingUser.ID = updatingAccount.User.ID

	if err := a.userSvc.Update(ctx, updatingUser); err != nil {
		return err
	}
	return nil
}

func (a *service) Authenticate(ctx context.Context, username, password string) (Account, error) {
	acct, err := a.Find(ctx, username)
	if err != nil {
		return NullAccount, ErrAccountNotFound
	}
	if err := a.secSvc.VerifyPasswordMatches(acct.Password, password); err != nil {
		return NullAccount, err
	}
	acct.LastLogin = time.Now().UTC()
	if err := a.accountRepo.Update(ctx, &acct); err != nil {
		return NullAccount, err
	}
	return acct, nil
}
