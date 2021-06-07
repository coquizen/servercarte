package ginHTTP

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/domain/authentication"
)

type authenticationMiddleware struct {
	authSvc authentication.Service
}

// NewMiddleWare is a constructor function to be used as an authentication middleware
func NewMiddleWare(authSvc authentication.Service) gin.HandlerFunc {
	return (&authenticationMiddleware{
		authSvc,
	}).handle
}

//
func (m *authenticationMiddleware) handle(ctx *gin.Context) {
	tokenString, err := m.authSvc.ExtractToken(ctx.Request)
	if err != nil {
		if err == authentication.ErrInvalidAccessToken {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	claims, err := m.authSvc.ParseTokenClaims(tokenString)
	if err != nil {
		if err == authentication.ErrInvalidAccessToken {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if err == authentication.ErrExpiredToken {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
 	}
	ctx.Set(authentication.CtxAuthenticationKey, claims)
}




