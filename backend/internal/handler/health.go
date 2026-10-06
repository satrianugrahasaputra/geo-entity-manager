package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse is returned by the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health reports that the process is up and able to serve requests.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
