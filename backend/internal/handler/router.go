package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouterConfig holds what the router needs to build the middleware chain.
type RouterConfig struct {
	Logger        *slog.Logger
	CORSOrigins   []string
	EntityHandler *EntityHandler
}

// NewRouter builds the Gin engine with middleware and routes.
func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	// Not behind a trusted proxy by default; ClientIP uses the socket address.
	_ = r.SetTrustedProxies(nil)

	r.Use(
		RequestID(),
		Logger(cfg.Logger),
		Recovery(cfg.Logger),
		CORS(cfg.CORSOrigins),
		BodyLimit(MaxBodyBytes),
	)

	// /healthz at the root for container health checks; also under /api/v1 so it
	// can be reached through the frontend's /api proxy.
	r.GET("/healthz", Health)
	api := r.Group("/api/v1")
	api.GET("/healthz", Health)
	
	if cfg.EntityHandler != nil {
		api.POST("/entities", cfg.EntityHandler.Create)
		api.GET("/entities/:id", cfg.EntityHandler.GetByID)
	}

	r.NoRoute(func(c *gin.Context) {
		abortWithError(c, http.StatusNotFound, CodeNotFound, "Resource tidak ditemukan", nil)
	})
	r.NoMethod(func(c *gin.Context) {
		abortWithError(c, http.StatusMethodNotAllowed, CodeMethodNotAllowed, "Method tidak diizinkan", nil)
	})

	return r
}
