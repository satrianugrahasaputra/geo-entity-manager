package handler

import (
	"time"

	"geo-entity-manager/backend/internal/model"
	"github.com/google/uuid"
)

// CreateEntityRequest is the payload for POST /entities
type CreateEntityRequest struct {
	Name        string             `json:"name" validate:"required,notblank,min=3,max=100"`
	Type        model.EntityType   `json:"type" validate:"required,oneof=vehicle iot_device facility other"`
	Status      model.EntityStatus `json:"status" validate:"required,oneof=active inactive maintenance offline"`
	Description string             `json:"description" validate:"max=500"`
	Latitude    float64            `json:"latitude" validate:"finite,min=-90,max=90"`
	Longitude   float64            `json:"longitude" validate:"finite,min=-180,max=180"`
}

// UpdateEntityRequest is the payload for PUT /entities/{id}
type UpdateEntityRequest struct {
	Name        string             `json:"name" validate:"required,notblank,min=3,max=100"`
	Type        model.EntityType   `json:"type" validate:"required,oneof=vehicle iot_device facility other"`
	Status      model.EntityStatus `json:"status" validate:"required,oneof=active inactive maintenance offline"`
	Description string             `json:"description" validate:"max=500"`
	Latitude    float64            `json:"latitude" validate:"finite,min=-90,max=90"`
	Longitude   float64            `json:"longitude" validate:"finite,min=-180,max=180"`
}

// PatchEntityRequest is the payload for PATCH /entities/{id}
type PatchEntityRequest struct {
	Name        *string             `json:"name,omitempty" validate:"omitempty,notblank,min=3,max=100"`
	Type        *model.EntityType   `json:"type,omitempty" validate:"omitempty,oneof=vehicle iot_device facility other"`
	Status      *model.EntityStatus `json:"status,omitempty" validate:"omitempty,oneof=active inactive maintenance offline"`
	Description *string             `json:"description,omitempty" validate:"omitempty,max=500"`
	Latitude    *float64            `json:"latitude,omitempty" validate:"omitempty,finite,min=-90,max=90"`
	Longitude   *float64            `json:"longitude,omitempty" validate:"omitempty,finite,min=-180,max=180"`
}

// ListEntityQuery represents query parameters for GET /entities.
type ListEntityQuery struct {
	Type   *model.EntityType   `form:"type" validate:"omitempty,oneof=vehicle iot_device facility other"`
	Status *model.EntityStatus `form:"status" validate:"omitempty,oneof=active inactive maintenance offline"`
	Search *string             `form:"search" validate:"omitempty,max=100"`
	Page   *int                `form:"page" validate:"omitempty,min=1"`
	Limit  *int                `form:"limit" validate:"omitempty,min=1,max=5000"`
}

// EntityResponse represents an entity sent to the client.
type EntityResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Type        model.EntityType   `json:"type"`
	Status      model.EntityStatus `json:"status"`
	Description string             `json:"description"`
	Latitude    float64            `json:"latitude"`
	Longitude   float64            `json:"longitude"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ListMeta holds pagination metadata.
type ListMeta struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// ListResponse is the envelope for GET /entities.
type ListResponse struct {
	Data []EntityResponse `json:"data"`
	Meta ListMeta         `json:"meta"`
}

// MapEntityToResponse converts a domain model to a response DTO.
func MapEntityToResponse(e *model.Entity) EntityResponse {
	return EntityResponse{
		ID:          e.ID,
		Name:        e.Name,
		Type:        e.Type,
		Status:      e.Status,
		Description: e.Description,
		Latitude:    e.Latitude,
		Longitude:   e.Longitude,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
