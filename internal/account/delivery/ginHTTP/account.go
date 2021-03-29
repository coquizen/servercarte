package ginHTTP

import (
	"fmt"
	"net/http"

	"github.com/CaninoDev/gastro/server/api/authentication"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/api/account"
	"github.com/CaninoDev/gastro/server/internal/logger"
)

type accountHandler struct {
	authSvc    authentication.Service
	accountSvc account.Service

}

// RegisterRoutes sets up account API endpoint using Gin.
func RegisterRoutes(authSvc authentication.Service, accountSvc account.Service, r *gin.Engine, authMiddleWare gin.HandlerFunc) {
	handler := accountHandler{authSvc, accountSvc}
	publicRoutes(handler, r)
	privateRoutes(handler, r, authMiddleWare)
}

func publicRoutes(handler accountHandler, router *gin.Engine) {
	router.POST("/login", handler.login)
}

func privateRoutes(handler accountHandler, router *gin.Engine, authMiddleWare gin.HandlerFunc) {

	routerGroup := router.Group("/accounts", authMiddleWare)
	routerGroup.GET("", handler.list)

	anotherRouterGroup := router.Group("/account", authMiddleWare)
	anotherRouterGroup.POST("", handler.create)
	anotherRouterGroup.PATCH("", handler.update)
	anotherRouterGroup.DELETE("", handler.delete)
}

func (h *accountHandler) create(ctx *gin.Context) {
	var newAccount account.NewAccountRequest
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err).SetMeta("malformed, partial, or missing registration request")
	}

	if err := h.accountSvc.New(ctx, newAccount); err != nil {
		if err := ctx.AbortWithError(http.StatusNotAcceptable, err).Error; err != nil {
			logger.Error.Println(err)
		}
		return
	}

	ctx.JSON(http.StatusOK, nil)
}


type credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *accountHandler) login(ctx *gin.Context) {
	var cred credentials
	if err := ctx.ShouldBindJSON(&cred); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	authenticationToken, err := h.accountSvc.Authenticate(ctx, cred.Username, cred.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err).SetMeta("unable to authenticate")
		return
	}

	ctx.JSON(http.StatusOK, authenticationToken)
}

func (h *accountHandler) list(ctx *gin.Context) {
		accounts, err := h.accountSvc.List(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err).SetMeta("unable to list accounts")
			return
		}
		ctx.JSON(http.StatusOK, accounts)
}

func (h *accountHandler) update(ctx *gin.Context) {
	var updateAccount account.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&updateAccount); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err).SetMeta("unable to process update request")
		return
	}
	accountID, exists := ctx.Get("accountID")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	updateAccount.ID = accountID.(uuid.UUID)
	if err := h.accountSvc.Update(ctx, updateAccount); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err).SetMeta("unable to update request")
		return
		}
	ctx.JSON(http.StatusOK, updateAccount)
}

type deleteRequest struct {
	password string
}

func (h *accountHandler) delete(ctx *gin.Context) {
	var deleteReq deleteRequest
	if err := ctx.ShouldBindJSON(&deleteReq); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	accountID, err := h.authSvc.ExtractClaims(ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unable to extract token: %v", err)).SetMeta("unable to extract token")
	}
	if err := h.accountSvc.Delete(ctx, accountID, deleteReq.password); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err).SetMeta("unable to delete request")
	}

	ctx.JSON(http.StatusOK, "account successfully deleted")
}

