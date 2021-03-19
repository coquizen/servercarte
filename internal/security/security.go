package security

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/helpers"
)

type BCrypt struct {
	Settings Settings
}

type Settings struct {
	Length        int
	MixedCase     bool
	AlphaNum      bool
	SpecialChar   bool
	CheckPrevious bool
}

func Bind(policy Settings) *BCrypt {
	return &BCrypt{Settings: policy}
}

func Initialize(cfg config.Security) *BCrypt {
	policy := Settings(cfg)
	return Bind(policy)
}

// ConfirmationChecker compares two given literal password inputs and returns whether they are equivalent.
func (s BCrypt) ConfirmationChecker(_ context.Context, password, confirmPassword string) bool {
	return password == confirmPassword
}

// VerifyPasswordMatches verifies that a given password is equivalent to its hashed form
func (s BCrypt) VerifyPasswordMatches(_ context.Context, hashedPW, givenPW string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(givenPW)) == nil
}

// Hash takes as input a given password and hashes it so it is suitable for persistent storage
func (s BCrypt) Hash(_ context.Context, password string) string {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encryptedPassword)
}

// IsValid validates whether a given string complies with the policy as declared by PasswordPolicy
func (s BCrypt) IsValid(_ context.Context, password string) error {
	// TODO: add conditional to check if password has been used previously.
	if s.Settings.Length > len(password) {
		return errors.New(fmt.Sprintf("password too short; should have %d characters", s.Settings.Length))
	}

	if s.Settings.MixedCase && !helpers.HasMixedCase(password) {
		return errors.New("password must have mixed lower and upper case letters")
	}
	if s.Settings.AlphaNum && !helpers.HasAlphaNum(password) {
		return errors.New("password must have both alphabet and number characters")
	}
	if s.Settings.SpecialChar && !helpers.HasSpecialChar(password) {
		return errors.New("password must have special characters")
	}
	return nil
}
