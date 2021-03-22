package gin

import (
	"fmt"
	"github.com/CaninoDev/gastro/server/api/account"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/internal/authentication"
)

type handler struct {
	authSvc authentication.Service
	svc     account.Service
}

// NewRoutes sets up menu API endpoint using Gin has the router.
func NewRoutes(svc account.Service, authSvc authentication.Service, r *gin.Engine) {
	h := handler{authSvc, svc}

	// public routes
	r.POST("/register", h.register)
	r.POST("/login", h.login)

	// private routes
	accountGroup := r.Group("/accounts", authSvc.Middleware())
	accountGroup.GET("", h.list)
	accountGroup.PATCH("", h.update)
	accountGroup.DELETE("", h.delete)

}


func (h *handler) register(ctx *gin.Context) {
	var newAccount account.newAccountRequest
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	if err := h.svc.New(ctx, newAccount); err != nil {
		ctx.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}


type credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *handler) login(ctx *gin.Context) {
	var cred credentials
	if err := ctx.ShouldBindJSON(&cred); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	authenticationToken, err := h.svc.Authenticate(ctx, cred.Username, cred.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
	}

	ctx.JSON(http.StatusOK, authenticationToken)
}

func (h *handler) list(ctx *gin.Context) {
	accounts, err := h.svc.List(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.JSON(http.StatusOK, accounts)
}

func (h *handler) update(ctx *gin.Context) {
	var updateAccount account.updateAccountRequest
	if err := ctx.ShouldBindJSON(&updateAccount); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	if err := h.authSvc.TokenValid(ctx.Request); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	claims, err := h.authSvc.ExtractTokenClaims(ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unable to extract token: %v", err))
	}
	updateAccount.ID = claims.AccountID

	if err := h.svc.Update(ctx, claims.AccountID, updateAccount); err != nil {
		if err := ctx.AbortWithError(http.StatusBadRequest, err).Error; err != nil {
			return
		}
	}
	ctx.JSON(http.StatusOK, updateAccount)
}

type deleteRequest struct {
	password string
}

func (h *handler) delete(ctx *gin.Context) {
	var deleteReq deleteRequest
	if err := ctx.ShouldBindJSON(&deleteReq); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	if err := h.authSvc.TokenValid(ctx.Request); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	claims, err := h.authSvc.ExtractTokenClaims(ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unable to extract token: %v", err))
	}
	if err := h.svc.Delete(ctx, claims.AccountID, deleteReq.password); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, "account successfully deleted")
}

