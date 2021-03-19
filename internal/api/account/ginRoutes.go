package account

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/internal/authentication"
	"github.com/CaninoDev/gastro/server/internal/model"
)

type handler struct {
	authSvc authentication.Service
	svc Service
}

// NewRoutes sets up menu API endpoint using Gin has the router.
func NewRoutes(svc Service, authSvc authentication.Service, r *gin.Engine) {
	h := handler{authSvc, svc}

	// public routes
	r.POST("/register", h.register)
	r.POST("/login", h.login)

	// private routes
	accountGroup := r.Group("", authSvc.Middleware())
	accountGroup.GET("/accounts", h.list)
	accountGroup.PATCH("/accounts/:id", h.update)
}

type createAccountRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Address1        string `json:"address_1" binding:"required"`
	Address2        string `json:"address_2,omitempty"`
	ZipCode         uint   `json:"zip_code" binding:"required"`
	Username        string `json:"username" binding:"required,min=3,alphanum"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
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

	tokenString, err := h.svc.Authenticate(ctx, cred.Username, cred.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, tokenString)


}

func (h *handler) register(ctx *gin.Context) {
	// // Handle the edge case where a user might currently be logged in.
	// tokenString := ctx.Request.Header.Get("Authorization")
	// if len(tokenString) > 0 {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "currently logged in existing account"})
	// }
	var newAccount createAccountRequest
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	if err := h.svc.Create(ctx, &newAccount); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{ "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (h *handler) list(ctx *gin.Context) {
	accounts, err := h.svc.List(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Address1        string `json:"address_1,omitempty"`
	Address2        string `json:"address_2,omitempty"`
	ZipCode         uint   `json:"zip_code,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"password_confirm,omitempty"`
	Email           string `json:"email,omitempty"`
	Role            model.AccessLevel `json:"role,omitempty"`
}

func (h *handler) update(ctx *gin.Context) {
	var updateAccount updateAccountRequest
	if err := ctx.ShouldBindJSON(&updateAccount); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	rawID := ctx.Param("id")

	h.svc.Update(ctx, rawID, &updateAccount)
}

