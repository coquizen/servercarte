package security

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/helpers"
)

type Security struct {
	passwordPolicy PasswordPolicy
}

type PasswordPolicy struct {
	Length        int
	MixedCase     bool
	AlphaNum      bool
	SpecialChar   bool
	CheckPrevious bool
}

func Bind(policy PasswordPolicy) *Security {
	return &Security{passwordPolicy: policy}
}
func Initialize(cfg config.Security) *Security {
	policy := PasswordPolicy(cfg)
	return Bind(policy)
}

func (s Security) Authenticate(_ context.Context, hashedPW, givenPW string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(givenPW)) == nil
}

func (s Security) Encrypt(_ context.Context, password string) string {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encryptedPassword)
}

func (s Security) IsValid(_ context.Context, password string) error {
	// TODO: add conditional to check if password has been used previously.
	if int(s.passwordPolicy.Length) > len(password) {
		return errors.New(fmt.Sprintf("password too short; should have %d characters", s.passwordPolicy.Length))
	}

	if s.passwordPolicy.MixedCase && !helpers.HasMixedCase(password) {
		return errors.New("password must have mixed lower and upper case letters")
	}
	if s.passwordPolicy.AlphaNum && !helpers.HasAlphaNum(password) {
		return errors.New("password must have both alphabet and number characters")
	}
	if s.passwordPolicy.SpecialChar && !helpers.HasSpecialChar(password) {
		return errors.New("password must have special characters")
	}
	return nil
}
