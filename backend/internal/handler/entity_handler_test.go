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
)

type fakeEntityService struct {
	createFunc  func(ctx context.Context, e *model.Entity) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*model.Entity, error)
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
