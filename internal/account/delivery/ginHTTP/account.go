package ginHTTP

import (
	"github.com/CaninoDev/gastro/server/domain/user"
	"net/http"

	"github.com/CaninoDev/gastro/server/domain/authentication"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/domain/account"
)

type accountHandler struct {
	authSvc    authentication.Service
	accountSvc account.Service
}

type accountRequest struct {
	Username string `json:"username"`
}

// RegisterRoutes sets up account API endpoint using Gin.
func RegisterRoutes(accountSvc account.Service, authSvc authentication.Service, r *gin.Engine, authMiddleWare gin.HandlerFunc, authorizationMiddleware gin.HandlerFunc) {
	handler := accountHandler{authSvc, accountSvc}
	publicRoutes(handler, r)
	privateRoutes(handler, r, authMiddleWare, authorizationMiddleware)
}

func publicRoutes(handler accountHandler, router *gin.Engine) {
	router.POST("/login", handler.login)
}

func privateRoutes(handler accountHandler, router *gin.Engine, authMiddleWare gin.HandlerFunc, authorizationMiddleware gin.HandlerFunc) {

	routerGroup := router.Group("/accounts", authMiddleWare, authorizationMiddleware)
	routerGroup.GET("", handler.list)

	anotherRouterGroup := router.Group("/account", authMiddleWare, authorizationMiddleware)
	anotherRouterGroup.GET("", handler.view)
	anotherRouterGroup.POST("", handler.create)
	anotherRouterGroup.PATCH("", handler.update)
	anotherRouterGroup.DELETE("", handler.delete)

}

func (h *accountHandler) create(ctx *gin.Context) {
	var newAccount account.NewAccountRequest
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	createdAccount, err := h.accountSvc.New(ctx, newAccount)
	if err != nil {
		if err == user.ErrUserAlreadyExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": createdAccount})
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

	acct, err := h.accountSvc.Authenticate(ctx, cred.Username, cred.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenString, err := h.authSvc.GenerateToken(ctx, acct.ID, acct.Username, int(acct.Role))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	claims, err := h.authSvc.ParseTokenClaims(tokenString)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString, "expiry": claims.Expiry})
}

func (h *accountHandler) list(ctx *gin.Context) {
	accounts, err := h.accountSvc.Accounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": accounts})
}

func (h *accountHandler) view(ctx *gin.Context) {
	var accountReq accountRequest
	if err := ctx.ShouldBindJSON(&accountReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	foundAccount, err := h.accountSvc.Find(ctx, accountReq.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": foundAccount})
}

func (h *accountHandler) update(ctx *gin.Context) {
	var updateAccount account.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&updateAccount); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	accountID, exists := ctx.Get("accountID")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	updateAccount.ID = accountID.(uuid.UUID)
	if err := h.accountSvc.Update(ctx, updateAccount); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "account successfully deleted")
}
