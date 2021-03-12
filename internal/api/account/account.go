package account

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/api/security"
	"github.com/CaninoDev/gastro/server/internal/api/user"
	"github.com/CaninoDev/gastro/server/internal/model"
)
// Account are the contracted methods to interact with GORM
type Account struct {
	accountRepo Repository
	userRepo user.Repository
	secSvc security.Security
}

func Bind(accountRepo Repository, userRepo user.Repository, secSvc security.Security) *Account {
	return &Account{
		accountRepo,userRepo, secSvc,
	}
}

func Initialize(accountRepo Repository, userRepo user.Repository, secSvc security.Security) *Account {
	return Bind(accountRepo, userRepo, secSvc)
}

func (a Account) Create(ctx context.Context, req *createRequest) error {
	var newUser model.User
	newUser.FirstName = req.FirstName
	newUser.LastName = req.LastName
	newUser.Addr = req.Addr
	newUser.ZipCode = req.ZipCode
	newUser.Email = req.Email
	if err := a.userRepo.View(ctx, &newUser); err == nil {
		return errors.New("account already exists")
	}

	var newAccount model.Account
	newAccount.Username = req.Username
	if err := a.accountRepo.Find(ctx, &newAccount); err == nil {
		return errors.New("username already exists")
	}
	if req.Password != req.PasswordConfirm {
		return errors.New("passwords don't match")
	}
	if err := a.secSvc.IsValid(ctx, req.Password); err != nil {
		return err
	}

	newAccount.Password = a.secSvc.Encrypt(ctx, req.Password)

	if err := a.userRepo.Create(ctx, &newUser); err != nil {
		return err
	}
	newAccount.UserID = newUser.ID
	if err := a.accountRepo.Create(ctx, &newAccount); err != nil {
		return err
	}
	return nil

}

func (a Account) ChangePassword(ctx context.Context, req passwordChangeRequest, userName string) error {
	var acct model.Account
	acct.Username = userName
	if req.NewPassword != req.NewPasswordConfirm {
		return errors.New("passwords don't match")
	}

	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return err
	}
	if ! a.secSvc.Authenticate(ctx, acct.Password, req.OldPassword) {
		return errors.New("password incorrect")
	}

	encryptedPW := a.secSvc.Encrypt(ctx, req.NewPassword)
	acct.Password = encryptedPW
	if err := a.accountRepo.Update(ctx, &acct); err != nil {
		return err
	}
	return nil
}

// Delete will delete the intended account
func (a Account) Delete(ctx context.Context, id string, passWord string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("malformed id")
	}
	var acct model.Account
	acct.ID = parsedID
	if err := a.accountRepo.Find(ctx, &acct); err != nil {
		return err
	}

	if ! a.secSvc.Authenticate(ctx, acct.Password, passWord) {
		return errors.New("password incorrect")
	}

	if err := a.accountRepo.Delete(ctx, &acct); err != nil {
		return err
	}
	return nil
}