package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"

	"geo-entity-manager/backend/internal/apperror"
	"geo-entity-manager/backend/internal/model"
)

// EntityRepository handles database operations for entities.
type EntityRepository struct {
	pool *pgxpool.Pool
}

// NewEntityRepository creates a new EntityRepository.
func NewEntityRepository(pool *pgxpool.Pool) *EntityRepository {
	return &EntityRepository{pool: pool}
}

// Create inserts a new entity into the database and sets the ID, CreatedAt, and UpdatedAt.
func (r *EntityRepository) Create(ctx context.Context, e *model.Entity) error {
	query := `
		INSERT INTO entities (name, type, status, description, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	
	err := r.pool.QueryRow(ctx, query,
		e.Name,
		e.Type,
		e.Status,
		e.Description,
		e.Latitude,
		e.Longitude,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create entity: %w", err)
	}

	return nil
}

// GetByID retrieves an entity by its ID.
func (r *EntityRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Entity, error) {
	query := `
		SELECT id, name, type, status, description, latitude, longitude, created_at, updated_at
		FROM entities
		WHERE id = $1
	`
	
	e := &model.Entity{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&e.ID,
		&e.Name,
		&e.Type,
		&e.Status,
		&e.Description,
		&e.Latitude,
		&e.Longitude,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperror.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get entity by id: %w", err)
	}

	return e, nil
}
