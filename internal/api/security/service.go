package security

import (
	"context"
)

type Service interface {
	ConfirmationChecker(ctx context.Context, password string, confirmPassword string) bool
	Authenticate(context.Context, string, string) bool
	Encrypt(context.Context, string) string
	IsValid(context.Context, string) error
}
