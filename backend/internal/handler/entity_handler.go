package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"geo-entity-manager/backend/internal/apperror"
	"geo-entity-manager/backend/internal/model"
)

type EntityService interface {
	Create(ctx context.Context, e *model.Entity) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Entity, error)
	List(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string, page, limit int) ([]model.Entity, int64, error)
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

// List handles GET /entities
func (h *EntityHandler) List(c *gin.Context) {
	var query ListEntityQuery
	
	// Bind valid query params
	if err := c.ShouldBindQuery(&query); err != nil {
		// Just a basic mapping for gin's binder errors; strict validation happens next
		HandleError(c, apperror.NewValidationError("Format query tidak valid", nil))
		return
	}

	// strict validation via custom DecodeQueryAndValidate
	if err := DecodeQueryAndValidate(c.Request, &query); err != nil {
		HandleError(c, err)
		return
	}

	page := 1
	if query.Page != nil {
		page = *query.Page
	}
	limit := 5000 // hard cap default
	if query.Limit != nil {
		limit = *query.Limit
	}

	entities, total, err := h.service.List(c.Request.Context(), query.Type, query.Status, query.Search, page, limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	var data []EntityResponse
	for _, e := range entities {
		data = append(data, MapEntityToResponse(&e))
	}
	// Always return empty array instead of null for empty results
	if data == nil {
		data = []EntityResponse{}
	}

	res := ListResponse{
		Data: data,
		Meta: ListMeta{
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}
	c.JSON(http.StatusOK, res)
}

