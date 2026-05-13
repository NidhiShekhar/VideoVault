package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	// Week 1 readiness: process is up and HTTP server is accepting requests.
	// In Week 2+, add Kafka/Postgres connectivity checks here.
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
