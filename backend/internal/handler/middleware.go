package handler

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// HeaderRequestID carries the request correlation ID.
	HeaderRequestID = "X-Request-ID"
	// MaxBodyBytes is the maximum accepted request body size (PRD §8: 1 MB).
	MaxBodyBytes int64 = 1 << 20

	ctxKeyRequestID  = "request_id"
	maxRequestIDSize = 64
)

// RequestID reuses a well-formed incoming X-Request-ID or generates a new one,
// stores it in the context, and echoes it in the response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(HeaderRequestID)
		if !validRequestID(id) {
			id = newRequestID()
		}
		c.Set(ctxKeyRequestID, id)
		c.Header(HeaderRequestID, id)
		c.Next()
	}
}

// RequestIDFrom returns the request ID stored by the RequestID middleware.
func RequestIDFrom(c *gin.Context) string {
	return c.GetString(ctxKeyRequestID)
}

func validRequestID(id string) bool {
	if id == "" || len(id) > maxRequestIDSize {
		return false
	}
	for _, r := range id {
		isAlnum := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
		if !isAlnum && r != '-' && r != '_' && r != '.' {
			return false
		}
	}
	return true
}

func newRequestID() string {
	b := make([]byte, 16)
	// crypto/rand.Read never returns an error on supported platforms (Go 1.24+ panics instead).
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Logger writes one structured log line per request.
func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		level := slog.LevelInfo
		switch {
		case status >= 500:
			level = slog.LevelError
		case status >= 400:
			level = slog.LevelWarn
		}
		log.LogAttrs(c.Request.Context(), level, "http request",
			slog.String("request_id", RequestIDFrom(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Duration("latency", time.Since(start)),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}

// Recovery turns panics into a 500 with the uniform error body, without leaking details.
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.ErrorContext(c.Request.Context(), "panic recovered",
					slog.String("request_id", RequestIDFrom(c)),
					slog.Any("panic", rec),
					slog.String("stack", string(debug.Stack())),
				)
				abortWithError(c, http.StatusInternalServerError, CodeInternal, "Terjadi kesalahan pada server", nil)
			}
		}()
		c.Next()
	}
}

// CORS allows cross-origin requests only from the configured origins.
// Requests from other origins get no CORS headers, so the browser blocks them.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := slices.Clone(allowedOrigins)
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && slices.Contains(allowed, origin) {
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Add("Vary", "Origin")
			h.Set("Access-Control-Expose-Headers", strings.Join([]string{"Location", HeaderRequestID}, ", "))

			if c.Request.Method == http.MethodOptions && c.GetHeader("Access-Control-Request-Method") != "" {
				h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				h.Set("Access-Control-Allow-Headers", "Content-Type, "+HeaderRequestID)
				h.Set("Access-Control-Max-Age", "600")
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}
		c.Next()
	}
}

// BodyLimit caps the request body size. Reading past the limit yields *http.MaxBytesError,
// which the JSON decoder maps to 413.
func BodyLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		}
		c.Next()
	}
}
