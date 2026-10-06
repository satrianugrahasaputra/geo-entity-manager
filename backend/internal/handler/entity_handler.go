package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"geo-entity-manager/backend/internal/model"
)

type EntityService interface {
	Create(ctx context.Context, e *model.Entity) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Entity, error)
}

type EntityHandler struct {
	service EntityService
}

func NewEntityHandler(service EntityService) *EntityHandler {
	return &EntityHandler{service: service}
}

// Create handles POST /entities
func (h *EntityHandler) Create(c *gin.Context) {
	var req CreateEntityRequest
	if err := DecodeAndValidate(c.Request, &req); err != nil {
		HandleError(c, err)
		return
	}

	entity := &model.Entity{
		Name:        req.Name,
		Type:        req.Type,
		Status:      req.Status,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	if err := h.service.Create(c.Request.Context(), entity); err != nil {
		HandleError(c, err)
		return
	}

	res := MapEntityToResponse(entity)
	c.Header("Location", fmt.Sprintf("/api/v1/entities/%s", res.ID.String()))
	c.JSON(http.StatusCreated, res)
}

// GetByID handles GET /entities/:id
func (h *EntityHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, CodeBadRequest, "Format UUID tidak valid", nil)
		return
	}

	entity, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, MapEntityToResponse(entity))
}
