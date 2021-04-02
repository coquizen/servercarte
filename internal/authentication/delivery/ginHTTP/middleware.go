package ginHTTP

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/authentication"
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
	claims, err := m.authSvc.ExtractClaims(ctx.Request)
	if err != nil {
		if err == authentication.ErrInvalidAccessToken {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Set(authentication.CtxAuthenticationKey, claims)
}




