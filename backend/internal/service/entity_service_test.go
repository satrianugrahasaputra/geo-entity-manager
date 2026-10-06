package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"geo-entity-manager/backend/internal/model"
)

type fakeRepository struct {
	createFunc  func(ctx context.Context, e *model.Entity) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*model.Entity, error)
}

func (r *fakeRepository) Create(ctx context.Context, e *model.Entity) error {
	if r.createFunc != nil {
		return r.createFunc(ctx, e)
	}
	return nil
}

func (r *fakeRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Entity, error) {
	if r.getByIDFunc != nil {
		return r.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func TestCreate(t *testing.T) {
	fakeRepo := &fakeRepository{}
	svc := NewEntityService(fakeRepo)

	t.Run("trims name", func(t *testing.T) {
		e := &model.Entity{Name: "  Pajero Sport  "}
		
		var savedName string
		fakeRepo.createFunc = func(ctx context.Context, entity *model.Entity) error {
			savedName = entity.Name
			return nil
		}

		err := svc.Create(context.Background(), e)
		require.NoError(t, err)
		assert.Equal(t, "Pajero Sport", savedName)
	})
}
