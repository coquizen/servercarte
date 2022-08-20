package bcrypto

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/coquizen/servercarte/internal/config"
	"github.com/coquizen/servercarte/internal/helpers"
)

var ErrPasswordsDoNotMatch = errors.New("passwords do not match")

type BCrypt struct {
	Length        int
	MixedCase     bool
	AlphaNum      bool
	SpecialChar   bool
	CheckPrevious bool
}

func NewSecurityFramework(cfg config.Security) *BCrypt {
	return &BCrypt{cfg.Length, cfg.MixedCase, cfg.AlphaNum, cfg.SpecialChar, cfg.CheckPrevious}
}

// ConfirmationChecker compares two given literal password inputs and returns whether they are equivalent.
func (s *BCrypt) ConfirmationChecker(password, confirmPassword string) error {
	if password != confirmPassword {
		return ErrPasswordsDoNotMatch
	}
	return nil
}

// VerifyPasswordMatches verifies that a given password is equivalent to its hashed form
func (s *BCrypt) VerifyPasswordMatches(hashedPW, givenPW string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(givenPW))
}

// Hash takes as input a given password and hashes it so it is suitable for persistent repository
func (s *BCrypt) Hash(password string) string {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encryptedPassword)
}

// IsValid validates whether a given string complies with the policy as declared by PasswordPolicy
func (s *BCrypt) IsValid(password string) error {
	// TODO: add conditional to check if password has been used previously.
	if s.Length > len(password) {
		return fmt.Errorf("password too short; should have %d characters", s.Length)
	}

	if s.MixedCase && !helpers.HasMixedCase(password) {
		return fmt.Errorf("password must have mixed lower and upper case letters")
	}
	if s.AlphaNum && !helpers.HasAlphaNum(password) {
		return fmt.Errorf("password must have both alphabet and number characters")
	}
	if s.SpecialChar && !helpers.HasSpecialChar(password) {
		return errors.New("password must have special characters")
	}
	return nil
}
