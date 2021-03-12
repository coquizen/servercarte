package account

import (

	"net/http"

	"github.com/gin-gonic/gin"

)

type ginHandler struct {
	svc Service
}

// NewGinRoutes sets up menu API endpoint using Gin has the router.
func NewGinRoutes(svc Service, r *gin.Engine) {
	h := ginHandler{svc}
	accountGroup := r.Group("")
	accountGroup.POST("/register", h.create)
	accountGroup.PATCH("/account/:id", h.changePassword)
	accountGroup.DELETE("/account/:id", h.delete)
}

type createRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Addr         string `json:"address" binding:"required"`
	ZipCode         uint    `json:"zip_code" binding:"required"`
	Username        string `json:"username" binding:"required,min=3,alphanum"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
}

// type updateRequest struct {
// 	FirstName       string `json:"first_name" binding:"required"`
// 	LastName        string `json:"last_name" binding:"required"`
// 	Address         string `json:"address" binding:"required"`
// 	ZipCode         uint    `json:"zip_code" binding:"required"`
// 	Email           string `json:"email" binding:"required,email"`
// }

type deleteRequest struct {
	Password string `json:"password"`
}

func (h *ginHandler) create(ctx *gin.Context) {
	var newAccount createRequest
	if err := ctx.ShouldBindJSON(&newAccount).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if newAccount.PasswordConfirm != newAccount.Password {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	if err := h.svc.Create(ctx, &newAccount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newAccount)
}

type passwordChangeRequest struct {
	OldPassword        string `json:"old_password" binding:"required,min=8"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"`
}

func (h *ginHandler) changePassword(ctx *gin.Context) {
	rawID := ctx.Param("id")

	var p passwordChangeRequest
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if p.NewPasswordConfirm != p.NewPasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	if err := h.svc.ChangePassword(ctx, p, rawID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully changed password"})
}

func (h *ginHandler) delete(ctx *gin.Context) {
	rawID := ctx.Param("id")
	var req deleteRequest
	if err := ctx.ShouldBindJSON(&req).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password is incorrect or empty"})
		return
	}

	if err := h.svc.Delete(ctx, rawID, req.Password); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted account"})

}



