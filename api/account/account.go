package account

import (
	"context"
	"errors"
	user2 "github.com/CaninoDev/gastro/server/api/user"
	"time"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/authentication"
	"github.com/CaninoDev/gastro/server/internal/model"
	"github.com/CaninoDev/gastro/server/internal/security"
)
// Account are the contracted methods to interact with GORM
type Account struct {
	accountRepo Repository
	userRepo    user2.Repository
	secSvc      security.Service
	authSvc     authentication.Service
}

func Bind(accountRepo Repository, userRepo user2.Repository, secSvc security.Service,
	authSvc authentication.Service) *Account {
	return &Account{
		accountRepo,userRepo, secSvc, authSvc,
	}
}

func Initialize(accountRepo Repository, userRepo user2.Repository, secSvc security.Service,
	authSvc authentication.Service) *Account {
	return Bind(accountRepo, userRepo, secSvc, authSvc)
}

func (a *Account) New(ctx context.Context, req NewAccountRequest) error {
	var newUser model.User
	newUser.FirstName = req.FirstName
	newUser.LastName = req.LastName
	newUser.Address1 = req.Address1
	newUser.Address2 = *req.Address2
	newUser.ZipCode = req.ZipCode
	newUser.Email = req.Email
	if err := a.userRepo.Search(ctx, &newUser); err == nil {
		return errors.New("account already exists")
	}

	if a.secSvc.ConfirmationChecker(ctx, req.Password, req.PasswordConfirm) == false {
		return errors.New("passwords don't match")
	}

	if err := a.secSvc.IsValid(ctx, req.Password); err != nil {
		return err
	}

	var newAccount model.Account
	newAccount.Username = req.Username
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

func (a *Account) Authenticate(ctx context.Context, username, password string) (string, error) {
	var acct model.Account
	acct.Username = username

	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return "", err
	}

	if !a.secSvc.VerifyPasswordMatches(ctx, acct.Password, password) {
		return "", errors.New("unauthorized")
	}

	token, err := a.authSvc.GenerateToken(ctx, &acct)
	if err != nil {
		return "", err
	}

	acct.LastLogin = time.Now()
	acct.Token = token

	if err := a.accountRepo.Update(ctx, &acct); err != nil {
		return "", err
	}
	return token, nil
}

func (a *Account) FindByUsername(ctx context.Context, username string) (*model.Account, error) {
	var acct model.Account
	acct.Username = username
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return &model.Account{}, err
	}
	return &acct, nil
}

func (a *Account) FindByToken(ctx context.Context, token string) (*model.Account, error) {
	var acct model.Account
	acct.Token = token
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return &model.Account{}, err
	}
	return &acct, nil
}

func (a *Account) ChangePassword(ctx context.Context, username, oldPassword, newPassword, confirmNewPassword string) error {
	var acct model.Account
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
func (a *Account) Delete(ctx context.Context, id uuid.UUID, passWord string) error {
	var acct model.Account
	acct.ID = id
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return err
	}

	if ! a.secSvc.VerifyPasswordMatches(ctx, acct.Password, passWord) {
		return errors.New("password incorrect")
	}

	if err := a.accountRepo.Delete(ctx, &acct); err != nil {
		return err
	}
	return nil
}

func (a *Account) RefreshAuthorization(ctx context.Context) error {
	claims := ctx.Value(authentication.AUTH_PROPS).(authentication.CustomClaims)
	var acct model.Account
	acct.Username = claims.Username
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return err
	}

	token, err := a.authSvc.GenerateToken(ctx, &acct)
	if err == nil {
		return err
	}
	acct.Token = token
	if err := a.accountRepo.Update(ctx, &acct); err != nil {
		return err
	}
	return errors.New("unauthorized; please re-login")
}

func (a *Account) List(ctx context.Context) (*[]model.Account, error) {
	var accounts []model.Account
	if err := a.accountRepo.All(ctx, &accounts); err != nil {
		return &accounts, err
	}
	return &accounts, nil
}

func (a *Account) Update(ctx context.Context, id uuid.UUID, request UpdateAccountRequest) error {
	var account model.Account
	account.ID = id
	if err := a.accountRepo.Find(ctx, &account); err != nil {
		return err
	}
	if err := a.accountRepo.Update(ctx, &account); err != nil {
		return err
	}
	var updateUser model.User
	updateUser.ID = account.UserID
	if err := a.userRepo.View(ctx, &updateUser); err != nil {
		return err
	}
	updateUser.Address1 = request.Address1
	updateUser.Address2 = *request.Address2
	updateUser.ZipCode = request.ZipCode
	updateUser.Email = request.Email
	if err := a.userRepo.Update(ctx, &updateUser); err != nil {
		return err
	}
	return nil
}