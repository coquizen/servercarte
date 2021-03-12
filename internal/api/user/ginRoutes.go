package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/model"
)

type ginHandler struct {
	svc Service
}

func NewGinRoutes(svc Service, r *gin.Engine) {
	h := ginHandler{svc}
	userGroup := r.Group("/user")
	userGroup.GET("/:id", h.view)
	userGroup.PATCH("/:id", h.update)
	userGroup.DELETE("/:id", h.delete)

}

// TODO: use the jwt token to unwrap claims of the currently logged in user
func (h *ginHandler) view(ctx *gin.Context) {
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
	FirstName string `json:"first_name,omitempty" gormDB:"unique,not null"`
	LastName  string `json:"last_name,omitempty" gormDB:"unique,not null"`
	Addr   string `json:"address" gormDB:"not null"`
	ZipCode   uint   `json:"zip_code" gormDB:"not null"`
}

func (h ginHandler) update(ctx *gin.Context) {
	var req = new(updateReq)
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	updateUser := &model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Addr:   req.Addr,
		ZipCode:   req.ZipCode,
	}
	updateUser.ID = parsedID
	if err := h.svc.Update(ctx, updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": updateUser})
}

func (h ginHandler) delete(ctx *gin.Context) {
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
