package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (a *JWT) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := a.parseToken(ctx.Request)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("no authorization header provided: %v", err))

		} else if isTokenExpired(token) {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token has expired: %v", err))

		} else if token.Valid {
			provisionClaimsToContext(ctx, token)

		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx.Next()
	}
}

func provisionClaimsToContext(ctx *gin.Context, token *jwt.Token) {
	ctxClaims := token.Claims.(*jwtWrappedClaims)
	ctx.Set("role", ctxClaims.Role)
	ctx.Set("accountID", ctxClaims.AccountID)
	ctx.Set("username", ctxClaims.Username)
}

func isTokenExpired(token *jwt.Token) bool {
	claims := token.Claims.(*jwtWrappedClaims)
	exp := time.Unix(claims.Expiry, 0).UTC()
	currTime := time.Now().UTC()
	return currTime.After(exp)
}
