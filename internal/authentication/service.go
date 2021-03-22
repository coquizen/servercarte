package authentication

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/model"
)


const AUTH_PROPS = "authProps"
type Token struct {
	Token string
	RefreshToken string
}

// CustomClaims are the custom Claims that identification authentication mechanism will certify.
type CustomClaims struct {
	AccountID uuid.UUID
	Username  string
	Email     string
	Role      model.AccessLevel
	Expiry    int64
}

// Service represents the minimum methods that the authentication system must implement
type Service interface {
	GenerateToken(ctx context.Context, acct *model.Account) (string, error)
	ExtractToken(req *http.Request) string
	ExtractTokenClaims(req *http.Request) (CustomClaims, error)
	TokenValid(req *http.Request) error
	Middleware() gin.HandlerFunc
}


