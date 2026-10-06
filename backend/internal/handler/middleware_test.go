package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newTestRouter(origins ...string) *gin.Engine {
	return NewRouter(RouterConfig{Logger: discardLogger(), CORSOrigins: origins})
}

func decodeError(t *testing.T, body io.Reader) ErrorResponse {
	t.Helper()
	var resp ErrorResponse
	require.NoError(t, json.NewDecoder(body).Decode(&resp))
	return resp
}

func TestHealth(t *testing.T) {
	for _, path := range []string{"/healthz", "/api/v1/healthz"} {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
		})
	}
}

func TestRequestID(t *testing.T) {
	tests := []struct {
		name       string
		incoming   string
		wantReused bool
	}{
		{name: "generated when missing", incoming: "", wantReused: false},
		{name: "reused when well formed", incoming: "abc-123_X.y", wantReused: true},
		{name: "replaced when too long", incoming: strings.Repeat("a", 65), wantReused: false},
		{name: "replaced when it has unsafe characters", incoming: "bad id\n<script>", wantReused: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			if tt.incoming != "" {
				req.Header.Set(HeaderRequestID, tt.incoming)
			}
			rec := httptest.NewRecorder()
			newTestRouter().ServeHTTP(rec, req)

			got := rec.Header().Get(HeaderRequestID)
			require.NotEmpty(t, got)
			if tt.wantReused {
				assert.Equal(t, tt.incoming, got)
			} else {
				assert.NotEqual(t, tt.incoming, got)
				assert.Len(t, got, 32)
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	r := gin.New()
	r.Use(RequestID(), Recovery(discardLogger()))
	r.GET("/boom", func(*gin.Context) { panic("secret internal detail") })

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/boom", nil))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "secret internal detail")
	resp := decodeError(t, rec.Body)
	assert.Equal(t, CodeInternal, resp.Error.Code)
	assert.NotNil(t, resp.Error.Details)
}

func TestCORS(t *testing.T) {
	const allowed = "http://localhost:5173"

	tests := []struct {
		name          string
		method        string
		origin        string
		preflight     bool
		wantStatus    int
		wantAllowOrig string
	}{
		{name: "allowed origin simple request", method: http.MethodGet, origin: allowed, wantStatus: http.StatusOK, wantAllowOrig: allowed},
		{name: "allowed origin preflight", method: http.MethodOptions, origin: allowed, preflight: true, wantStatus: http.StatusNoContent, wantAllowOrig: allowed},
		{name: "disallowed origin gets no CORS headers", method: http.MethodGet, origin: "http://evil.test", wantStatus: http.StatusOK, wantAllowOrig: ""},
		{name: "no origin header", method: http.MethodGet, wantStatus: http.StatusOK, wantAllowOrig: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/healthz", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.preflight {
				req.Header.Set("Access-Control-Request-Method", http.MethodPost)
			}
			rec := httptest.NewRecorder()
			newTestRouter(allowed).ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantAllowOrig, rec.Header().Get("Access-Control-Allow-Origin"))
			if tt.preflight {
				assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), http.MethodPatch)
				assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
			}
		})
	}
}

func TestBodyLimit(t *testing.T) {
	tests := []struct {
		name        string
		size        int64
		wantTooLong bool
	}{
		{name: "exactly at limit", size: MaxBodyBytes, wantTooLong: false},
		{name: "over limit", size: MaxBodyBytes + 1, wantTooLong: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var readErr error
			r := gin.New()
			r.Use(BodyLimit(MaxBodyBytes))
			r.POST("/echo", func(c *gin.Context) {
				_, readErr = io.ReadAll(c.Request.Body)
				c.Status(http.StatusOK)
			})

			body := strings.NewReader(strings.Repeat("x", int(tt.size)))
			r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/echo", body))

			var maxErr *http.MaxBytesError
			assert.Equal(t, tt.wantTooLong, errors.As(readErr, &maxErr))
		})
	}
}

func TestFallbackRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantCode   string
	}{
		{name: "unknown route", method: http.MethodGet, path: "/nope", wantStatus: http.StatusNotFound, wantCode: CodeNotFound},
		{name: "wrong method", method: http.MethodDelete, path: "/healthz", wantStatus: http.StatusMethodNotAllowed, wantCode: CodeMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestRouter().ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantCode, decodeError(t, rec.Body).Error.Code)
		})
	}
}
