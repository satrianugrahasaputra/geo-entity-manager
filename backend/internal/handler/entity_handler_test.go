package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"geo-entity-manager/backend/internal/apperror"
	"geo-entity-manager/backend/internal/model"
	"geo-entity-manager/backend/internal/service"
)

type fakeEntityService struct {
	createFunc  func(ctx context.Context, e *model.Entity) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*model.Entity, error)
	listFunc    func(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string, page, limit int) ([]model.Entity, int64, error)
}

func (s *fakeEntityService) Create(ctx context.Context, e *model.Entity) error {
	if s.createFunc != nil {
		return s.createFunc(ctx, e)
	}
	e.ID = uuid.New()
	return nil
}

func (s *fakeEntityService) GetByID(ctx context.Context, id uuid.UUID) (*model.Entity, error) {
	if s.getByIDFunc != nil {
		return s.getByIDFunc(ctx, id)
	}
	return &model.Entity{ID: id, Name: "Test"}, nil
}

func (s *fakeEntityService) List(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string, page, limit int) ([]model.Entity, int64, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, typeFilter, statusFilter, search, page, limit)
	}
	return nil, 0, nil
}

func (s *fakeEntityService) Update(ctx context.Context, id uuid.UUID, e *model.Entity) (*model.Entity, error) {
	e.ID = id
	return e, nil
}

func (s *fakeEntityService) Patch(ctx context.Context, id uuid.UUID, req service.PatchRequest) (*model.Entity, error) {
	return &model.Entity{ID: id, Name: "Patched"}, nil
}

func TestCreateEntity(t *testing.T) {
	fakeSvc := &fakeEntityService{}
	router := NewRouter(RouterConfig{Logger: discardLogger(), EntityHandler: NewEntityHandler(fakeSvc)})

	t.Run("success", func(t *testing.T) {
		body := `{"name":"Truk 01","type":"vehicle","status":"active","latitude":-7.1,"longitude":110.4}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/entities", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NotEmpty(t, rec.Header().Get("Location"))

		var resp EntityResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		assert.Equal(t, "Truk 01", resp.Name)
	})

	t.Run("validation error", func(t *testing.T) {
		body := `{"name":"","type":"vehicle","status":"active","latitude":-7.1,"longitude":110.4}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/entities", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		resp := decodeError(t, rec.Body)
		assert.Equal(t, CodeValidation, resp.Error.Code)
	})
}

func TestGetByID(t *testing.T) {
	fakeSvc := &fakeEntityService{}
	router := NewRouter(RouterConfig{Logger: discardLogger(), EntityHandler: NewEntityHandler(fakeSvc)})
	validID := uuid.New()

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/entities/"+validID.String(), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/entities/invalid-uuid", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		resp := decodeError(t, rec.Body)
		assert.Equal(t, CodeBadRequest, resp.Error.Code)
	})

	t.Run("not found", func(t *testing.T) {
		fakeSvc.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*model.Entity, error) {
			return nil, apperror.ErrNotFound
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/entities/"+validID.String(), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		resp := decodeError(t, rec.Body)
		assert.Equal(t, CodeNotFound, resp.Error.Code)
	})
}

func TestListEntities(t *testing.T) {
	fakeSvc := &fakeEntityService{}
	router := NewRouter(RouterConfig{Logger: discardLogger(), EntityHandler: NewEntityHandler(fakeSvc)})

	t.Run("success", func(t *testing.T) {
		fakeSvc.listFunc = func(ctx context.Context, typeFilter *model.EntityType, statusFilter *model.EntityStatus, search *string, page int, limit int) ([]model.Entity, int64, error) {
			return []model.Entity{{ID: uuid.New(), Name: "E1"}}, 1, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/entities?page=1&limit=10", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp ListResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		assert.Equal(t, int64(1), resp.Meta.Total)
		assert.Len(t, resp.Data, 1)
	})

	t.Run("invalid query param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/entities?unknown=param", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Expecting 422 Unprocessable Entity due to strict query validation
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		resp := decodeError(t, rec.Body)
		assert.Equal(t, CodeValidation, resp.Error.Code)
	})
}

