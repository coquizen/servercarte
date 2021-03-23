package ginHTTP

import (
	"log"
	"net/http"

	"github.com/CaninoDev/gastro/server/api"
	authentication2 "github.com/CaninoDev/gastro/server/api/authentication"
	"github.com/CaninoDev/gastro/server/api/menu"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type menuHandler struct {
	svc     menu.Service
	authSvc authentication2.Service
}

// NewMenuRoutes sets up menu API endpoint using Gin has the router.
func NewMenuRoutes(svc menu.Service, authSvc authentication2.Service, r *gin.Engine) {
	h := menuHandler{svc, authSvc}
	menuGroup := r.Group("/api/v1")
	menuViewGroup := menuGroup.Group("")
	menuViewGroup.GET("/sections", h.listSections)
	menuViewGroup.GET("/sections/:id", h.findSectionByID)
	menuViewGroup.GET("/items", h.listItems)
	menuViewGroup.GET("/items/:id", h.findItemByID)

	menuEditGroup := menuGroup.Group("")
	menuEditGroup.POST("/sections", h.createSection)
	menuEditGroup.PATCH("/sections/:id", h.updateSection)
	menuEditGroup.DELETE("/sections/:id", h.deleteSection)
	menuEditGroup.POST("/items", h.createItem)
	menuEditGroup.PATCH("/items/:id", h.updateItem)
	menuEditGroup.DELETE("/items/:id", h.deleteItem)
}

// --- Sections --- //
func (h *menuHandler) listSections(ctx *gin.Context) {
	sections, err := h.svc.Sections(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sections})
}

func (h *menuHandler) findSectionByID(ctx *gin.Context) {
	rawID := ctx.Param("id")
	log.Printf("ID: %s", rawID)
	section, err := h.svc.SectionByID(ctx, rawID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": section})
}

// createSection creates a new section.
func (h *menuHandler) createSection(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != api.Admin {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var section api.Section

	if err := ctx.ShouldBindJSON(&section); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.NewSection(ctx, &section); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": section})
}

// updateSection update section's data.
func (h *menuHandler) updateSection(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != api.Admin {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var section api.Section

	if err := ctx.ShouldBindJSON(&section); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawID := ctx.Param("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section.ID = id
	if err := h.svc.UpdateSectionData(ctx, &section); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": section})
}

func (h *menuHandler) deleteSection(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != api.Admin {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	rawID := ctx.Param("id")
	if err := h.svc.DeleteSection(ctx, rawID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "section deleted"})
}

// ---  Item  --- //
func (h *menuHandler) listItems(ctx *gin.Context) {
	items, err := h.svc.Items(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

// createSection creates a new section.
func (h *menuHandler) createItem(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != api.Admin {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var item api.Item

	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.NewItem(ctx, &item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

// updateSection creates a new section.
func (h *menuHandler) updateItem(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != api.Admin {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	rawID := ctx.Param("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var item api.Item

	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = id
	if err := h.svc.UpdateItemData(ctx, &item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}
func (h *menuHandler) findItemByID(ctx *gin.Context) {
	rawID := ctx.Param("id")
	item, err := h.svc.ItemByID(ctx, rawID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *menuHandler) deleteItem(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	if !exists || role != api.Admin {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	rawID := ctx.Param("id")

	if err := h.svc.DeleteItem(ctx, rawID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "item deleted"})
}
