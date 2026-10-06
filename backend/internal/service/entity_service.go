package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"geo-entity-manager/backend/internal/model"
)

// EntityRepository defines the expected database operations.
type EntityRepository interface {
	Create(ctx context.Context, e *model.Entity) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Entity, error)
	List(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string, limit, offset int) ([]model.Entity, error)
	Count(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string) (int64, error)
	Update(ctx context.Context, e *model.Entity) error
}

// EntityService contains the business logic for entities.
type EntityService struct {
	repo EntityRepository
}

// NewEntityService creates a new EntityService.
func NewEntityService(repo EntityRepository) *EntityService {
	return &EntityService{repo: repo}
}

// Create processes and creates a new entity.
func (s *EntityService) Create(ctx context.Context, e *model.Entity) error {
	e.Name = strings.TrimSpace(e.Name)
	return s.repo.Create(ctx, e)
}

// GetByID retrieves an entity.
func (s *EntityService) GetByID(ctx context.Context, id uuid.UUID) (*model.Entity, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves paginated entities and the total count.
func (s *EntityService) List(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string, page, limit int) ([]model.Entity, int64, error) {
	offset := (page - 1) * limit
	count, err := s.repo.Count(ctx, typeFilter, statusFilter, search)
	if err != nil {
		return nil, 0, err
	}

	entities, err := s.repo.List(ctx, typeFilter, statusFilter, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return entities, count, nil
}

// Update completely replaces an entity's fields.
func (s *EntityService) Update(ctx context.Context, id uuid.UUID, e *model.Entity) (*model.Entity, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	e.ID = existing.ID
	e.Name = strings.TrimSpace(e.Name)
	if err := s.repo.Update(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

// PatchRequest holds the partial update fields.
type PatchRequest struct {
	Name        *string
	Type        *model.EntityType
	Status      *model.EntityStatus
	Description *string
	Latitude    *float64
	Longitude   *float64
}

// Patch partially updates an entity.
func (s *EntityService) Patch(ctx context.Context, id uuid.UUID, req PatchRequest) (*model.Entity, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = strings.TrimSpace(*req.Name)
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Latitude != nil {
		existing.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		existing.Longitude = *req.Longitude
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
