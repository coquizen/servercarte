package security

import (
	"context"
)

type PasswordService interface {
	Authenticate(context.Context, string, string) (bool, error)
	IsValid(context.Context, string) error
	Validate(context.Context, string, string) bool
}

type AuthenticationService interface {
	Encrypt(context.Context, string) string
	Token(context.Context, string) (string, error)
}
