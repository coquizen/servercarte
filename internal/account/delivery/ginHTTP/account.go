package ginHTTP

import (
	"net/http"

	"github.com/CaninoDev/gastro/server/authentication"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/domain/account"
	"github.com/CaninoDev/gastro/server/internal/logger"
)

type accountHandler struct {
	authSvc    authentication.Service
	accountSvc account.Service

}

// RegisterRoutes sets up account API endpoint using Gin.
func RegisterRoutes(authSvc authentication.Service, accountSvc account.Service, r *gin.Engine, authMiddleWare gin.HandlerFunc, authorizationMiddleware gin.HandlerFunc) {
	handler := accountHandler{authSvc, accountSvc}
	publicRoutes(handler, r)
	privateRoutes(handler, r, authMiddleWare,authorizationMiddleware)
}

func publicRoutes(handler accountHandler, router *gin.Engine) {
	router.POST("/login", handler.login)
}

func privateRoutes(handler accountHandler, router *gin.Engine, authMiddleWare gin.HandlerFunc, authorizationMiddleware gin.HandlerFunc) {

	routerGroup := router.Group("/accounts", authMiddleWare,authorizationMiddleware)
	routerGroup.GET("", handler.list)

	anotherRouterGroup := router.Group("/account", authMiddleWare, authorizationMiddleware)
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

	account, err := h.accountSvc.Authenticate(ctx, cred.Username, cred.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err).SetMeta("unable to authenticate")
		return
	}
	tokenString, err := h.authSvc.GenerateToken(ctx, account.ID, account.Username, int(account.Role))

	ctx.JSON(http.StatusOK, tokenString)
}

func (h *accountHandler) list(ctx *gin.Context) {
		accounts, err := h.accountSvc.Accounts(ctx)
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


func (h *accountHandler) delete(ctx *gin.Context) {
	rawID := ctx.Param("id")
	delID, err := uuid.Parse(rawID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := h.accountSvc.Delete(ctx, delID); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err).SetMeta("unable to delete request")
		return
	}

	ctx.JSON(http.StatusOK, "account successfully deleted")
}

