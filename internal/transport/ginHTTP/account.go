package ginHTTP

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/api/account"
	"github.com/CaninoDev/gastro/server/internal/logger"
	"github.com/CaninoDev/gastro/server/internal/model"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/CaninoDev/gastro/server/internal/authentication"
)

type accountHandler struct {
	authSvc authentication.Service
	svc     account.Service
}

// NewAccountRoutes sets up menu API endpoint using Gin has the router.
func NewAccountRoutes(svc account.Service, authSvc authentication.Service, r *gin.Engine) {
	h := accountHandler{authSvc, svc}

	// public routes
	r.POST("/register", h.register)
	r.POST("/login", h.login)

	// private routes
	accountGroup := r.Group("/accounts", authSvc.Middleware())
	accountGroup.GET("", h.list)
	accountGroup.PATCH("", h.update)
	accountGroup.DELETE("", h.delete)

}


func (h *accountHandler) register(ctx *gin.Context) {
	var newAccount account.NewAccountRequest
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		if err := ctx.AbortWithError(http.StatusUnprocessableEntity, err).Error; err != nil {
			logger.Error.Println(err)
		}
		return
	}

	if err := h.svc.New(ctx, newAccount); err != nil {
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

	// loginAcct, err := h.svc.FindByUsername(ctx, cred.Username)
	// if err != nil {
	// 	ctx.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }
	authenticationToken, err := h.svc.Authenticate(ctx, cred.Username, cred.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusOK, authenticationToken)
}

func (h *accountHandler) list(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	logger.Info.Println(role)
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if role == model.Admin {
		accounts, err := h.svc.List(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": accounts})
		return
	}
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

func (h *accountHandler) update(ctx *gin.Context) {
	var updateAccount account.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&updateAccount); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	accountID, exists := ctx.Get("accountID")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	updateAccount.ID = accountID.(uuid.UUID)
	if err := h.svc.Update(ctx, updateAccount); err != nil {
		if err := ctx.AbortWithError(http.StatusBadRequest, err).Error; err != nil {
			return
		}
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

