package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(uploadHandler *UploadHandler, healthHandler *HealthHandler) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/healthz", healthHandler.Healthz)
	r.GET("/readyz", healthHandler.Readyz)

	videos := r.Group("/api/v1/videos")
	videos.POST("/upload", uploadHandler.UploadVideo)
	videos.GET("/upload/status/:videoId", uploadHandler.GetUploadStatus)

	return r
}
