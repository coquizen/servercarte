package account

import (
	"context"
	"errors"

	"github.com/CaninoDev/gastro/server/api/authentication"
	"github.com/CaninoDev/gastro/server/api/security"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/api/user"
)

// service are the contracted methods to interact with GORM
type service struct {
	accountRepo Repository
	userRepo    user.Repository
	secSvc      security.Service
	authSvc     authentication.Service
}

func NewService(accountRepo Repository, userRepo user.Repository, secSvc security.Service,
	authSvc authentication.Service) *service {
	return &service{accountRepo, userRepo, secSvc, authSvc}
}

func (a *service) New(ctx context.Context, req NewAccountRequest) error {
	var newUser user.User
	newUser.FirstName = req.FirstName
	newUser.LastName = req.LastName
	newUser.Address1 = req.Address1
	newUser.Address2 = *req.Address2
	newUser.ZipCode = req.ZipCode
	newUser.Email = req.Email
	if err := a.userRepo.Search(ctx, &newUser); err == nil {
		return errors.New("account already exists")
	}

	if !a.secSvc.ConfirmationChecker(ctx, req.Password, req.PasswordConfirm) {
		return errors.New("passwords don't match")
	}

	if err := a.secSvc.IsValid(ctx, req.Password); err != nil {
		return err
	}

	var newAccount Account
	newAccount.Username = req.Username
	newAccount.Role = 0
	if err := a.accountRepo.Find(ctx, &newAccount); err == nil {
		return errors.New("username already exists")
	}

	newAccount.Password = a.secSvc.Hash(ctx, req.Password)

	if err := a.userRepo.Create(ctx, &newUser); err != nil {
		return err
	}
	newAccount.UserID = newUser.ID

	if err := a.accountRepo.Create(ctx, &newAccount); err != nil {
		return err
	}
	return nil
}

func (a *service) Authenticate(ctx context.Context, username, password string) (string, error) {
	acct, err := a.Find(ctx, username)
	if err != nil {
		return "", err
	}
	if !a.secSvc.VerifyPasswordMatches(ctx, acct.Password, password) {
		return "", errors.New("password do not match")
	}
	return a.authSvc.GenerateToken(ctx, acct.ID )
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
func (a *service) Delete(ctx context.Context, id uuid.UUID, password string) error {
	var acct Account
	acct.ID = id
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return err
	}

	if !a.secSvc.VerifyPasswordMatches(ctx, acct.Password, password) {
		return errors.New("password incorrect")
	}

	if err := a.accountRepo.Delete(ctx, &acct); err != nil {
		return err
	}
	return nil
}

func (a *service) RefreshAuthorization(ctx context.Context) error {
	userN := ctx.Value("username")
	var acct Account
	acct.Username = userN.(string)
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return err
	}

	token, err := a.authSvc.GenerateToken(ctx, acct.ID)
	if err == nil {
		return err
	}
	acct.Token = token
	if err := a.accountRepo.Update(ctx, &acct); err != nil {
		return err
	}
	return errors.New("unauthorized; please re-login")
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
	if err := a.accountRepo.Find(ctx, &account); err != nil {
		return err
	}
	if err := a.accountRepo.Update(ctx, &account); err != nil {
		return err
	}
	var updateUser user.User
	updateUser.ID = account.UserID
	if err := a.userRepo.View(ctx, &updateUser); err != nil {
		return err
	}
	updateUser.Address1 = request.Address1
	if request.Address2 != nil {
		updateUser.Address2 = *request.Address2
	}
	updateUser.ZipCode = request.ZipCode
	updateUser.Email = request.Email
	if err := a.userRepo.Update(ctx, &updateUser); err != nil {
		return err
	}
	return nil
}
