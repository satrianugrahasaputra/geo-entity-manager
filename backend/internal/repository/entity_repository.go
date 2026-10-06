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

// buildListQuery constructs the WHERE clause and arguments.
func buildListQuery(typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argID := 1

	if typeFilter != nil && *typeFilter != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argID))
		args = append(args, *typeFilter)
		argID++
	}
	if statusFilter != nil && *statusFilter != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
		args = append(args, *statusFilter)
		argID++
	}
	if search != nil && *search != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argID))
		args = append(args, "%"+*search+"%")
		argID++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			where += " AND " + conditions[i]
		}
	}
	return where, args
}

// Count returns the total number of matching entities.
func (r *EntityRepository) Count(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string) (int64, error) {
	where, args := buildListQuery(typeFilter, statusFilter, search)
	query := "SELECT COUNT(*) FROM entities " + where

	var count int64
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count entities: %w", err)
	}
	return count, nil
}

// List returns a paginated list of matching entities.
func (r *EntityRepository) List(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string, limit, offset int) ([]model.Entity, error) {
	where, args := buildListQuery(typeFilter, statusFilter, search)
	argID := len(args) + 1
	
	query := fmt.Sprintf(`
		SELECT id, name, type, status, description, latitude, longitude, created_at, updated_at
		FROM entities
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argID, argID+1)
	
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query entities: %w", err)
	}
	defer rows.Close()

	var results []model.Entity
	for rows.Next() {
		var e model.Entity
		if err := rows.Scan(
			&e.ID, &e.Name, &e.Type, &e.Status, &e.Description,
			&e.Latitude, &e.Longitude, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan entity: %w", err)
		}
		results = append(results, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

// Update modifies an existing entity.
func (r *EntityRepository) Update(ctx context.Context, e *model.Entity) error {
	query := `
		UPDATE entities
		SET name = $1, type = $2, status = $3, description = $4, latitude = $5, longitude = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		e.Name, e.Type, e.Status, e.Description, e.Latitude, e.Longitude, e.ID,
	).Scan(&e.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update entity: %w", err)
	}

	return nil
}

