package ginHTTP

import (
	"net/http"

	"github.com/CaninoDev/gastro/server/api"
	"github.com/CaninoDev/gastro/server/api/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandler struct {
	svc user.Service
}

// NewUserRoutes setups the api endpoint for the user
func NewUserRoutes(svc user.Service, r *gin.Engine) {
	h := userHandler{svc}
	userGroup := r.Group("/user")
	userGroup.GET("/:id", h.view)
	userGroup.PATCH("/:id", h.update)
	userGroup.DELETE("/:id", h.delete)

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
	return
}

type updateReq struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Address1   string `json:"address_1"`
	Address2  string `json:"address_2"`
	ZipCode   uint   `json:"zip_code"`
}

func (h userHandler) update(ctx *gin.Context) {
	var req updateReq
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	updateUser := &api.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Address1:  req.Address1,
		Address2: req.Address2,
		ZipCode:   req.ZipCode,
	}
	updateUser.ID = parsedID
	if err := h.svc.Update(ctx, updateUser); err != nil {
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
