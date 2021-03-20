package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/internal/authentication"
)



func (a *JWT) Middleware() gin.HandlerFunc {
		return func(ctx *gin.Context) {
			token, err := a.parseToken(ctx.Request)
			if err != nil {
				ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("No authorization header provided: %v", err))
			}

			if isTokenExpired(*token) == true {
				ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Token has expired: %v", err))
			}
			if token.Valid {
				provisionClaimsToContext(ctx, token)
				ctx.Next()
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			ctx.Next()
		}
}

func provisionClaimsToContext(ctx *gin.Context, token *jwt.Token) {
	ctxClaims := token.Claims.(jwtWrappedClaims).CustomClaims
	ctx.Set(authentication.AUTH_PROPS, ctxClaims)
}

func isTokenExpired(token jwt.Token) bool {
	claims := token.Claims.(jwtWrappedClaims)

	exp, err := time.Parse(time.RFC3339, string(claims.Expiry))
	if err != nil {
		return true
	}
	currTime := time.Now().UTC()
	if currTime.After(exp) == true {
		return true
	}
	return false
}
