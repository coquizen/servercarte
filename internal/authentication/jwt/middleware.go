package jwt

import (
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
			http.Error(ctx.Writer, "unable to parse token", http.StatusUnauthorized)
			return
		}
		if isTokenExpired(*token) == true {
			http.Error(ctx.Writer, "token has expired", http.StatusUnauthorized)
			return
		}
		provisionClaimsToContext(ctx, token)
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
