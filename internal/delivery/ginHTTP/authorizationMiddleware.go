package ginHTTP

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/domain/account"
	"github.com/CaninoDev/gastro/server/domain/authentication"
)

func AuthorizationMiddleware(accessLevel account.AccessLevel) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get(authentication.CtxAuthenticationKey)
		if !exists {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if account.AccessLevel(claims.(authentication.CustomClaims).Role) != accessLevel {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}