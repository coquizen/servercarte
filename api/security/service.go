package security

import (
	"context"
)

type Service interface {
	ConfirmationChecker(ctx context.Context, password string, confirmPassword string) bool
	VerifyPasswordMatches(ctx context.Context, hashedPW string, password string) bool
	Hash(ctx context.Context, password string) string
	IsValid(ctx context.Context, password string) error
}

type Security struct {
	Service
}

func NewService(secSvc Service) *Security {
	return &Security{secSvc}
}

