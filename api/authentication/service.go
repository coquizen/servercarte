package authentication

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const CtxAuthenticationKey = "auth"

// CustomClaims are the custom Claims that identification authentication mechanism will certify.
type CustomClaims struct {
	AccountID uuid.UUID
	Expiry    int64
}

// Service represents the minimum methods that the authentication system must implement
type Service interface {
	GenerateToken(ctx context.Context, accountID uuid.UUID) (string, error)
	ExtractClaims(req *http.Request) (uuid.UUID, error)
	TokenValid(tokenString string) error
}

type authentication struct {
	Service
}

// NewService returns an authentication - compliant service instance. Must satisfy the Service interface.
func NewService(authSvc  Service) *authentication {
	return &authentication{authSvc}
}


