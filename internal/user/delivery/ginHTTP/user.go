package ginHTTP

import (
	"net/http"
	"strings"

	"github.com/CaninoDev/gastro/server/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandler struct {
	svc user.Service
}

// RegisterRoutes setups the api endpoint for the user
func RegisterRoutes(svc user.Service, r *gin.Engine) {
	h := userHandler{svc}
	userGroup := r.Group("/user")
	userGroup.GET("/:id", h.view)
	userGroup.PATCH("/:id", h.update)
	userGroup.DELETE("/:id", h.delete)
	r.GET("/users", h.list)
}

func (h *userHandler) list(ctx *gin.Context) {
	users, err := h.svc.Users(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

// TODO: use the jwt token to unwrap claims of the currently logged in user
func (h *userHandler) view(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	userInformation, err := h.svc.View(ctx, parsedID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": userInformation})
}

type updateRequest struct {
	Address1        *string `json:"address_1,omitempty"`
	Address2        *string `json:"address_2,omitempty"`
	ZipCode         *uint   `json:"zip_code,omitempty"`
	TelephoneNumber *string `json:"telephone_number,omitempty"`
	Email           *string `json:"email,omitempty"`
}

func (h userHandler) update(ctx *gin.Context) {
	var req updateRequest
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var updateUser user.User
	if req.Address1 != nil {
		updateUser.Address1 = *req.Address1
	}
	if req.Address2 != nil {
		updateUser.Address2 = req.Address2
	}
	if req.ZipCode != nil {
		updateUser.ZipCode = *req.ZipCode
	}
	if req.TelephoneNumber != nil {
		updateUser.TelephoneNumber = trimPhoneString(*req.TelephoneNumber)
	}

	updateUser.ID = parsedID
	if err := h.svc.Update(ctx, &updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": updateUser})
}

func (h userHandler) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := h.svc.Delete(ctx, parsedID); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

func trimPhoneString(str string) string {
	return strings.Trim(str, "()- ")
}
