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
