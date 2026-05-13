package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/nidhi/video-vault/services/upload-service/internal/service"
)

type UploadHandler struct {
	videoService *service.VideoService
}

func NewUploadHandler(videoService *service.VideoService) *UploadHandler {
	return &UploadHandler{videoService: videoService}
}

type uploadRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
}

func (h *UploadHandler) UploadVideo(c *gin.Context) {
	var req uploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	video, err := h.videoService.Upload(c.Request.Context(), service.UploadVideoInput{
		Title:       req.Title,
		Description: req.Description,
		Filename:    req.Filename,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"videoId": video.ID,
		"status":  video.Status,
		"message": "upload accepted",
	})
}

func (h *UploadHandler) GetUploadStatus(c *gin.Context) {
	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "videoId is required"})
		return
	}

	video, err := h.videoService.GetStatus(c.Request.Context(), videoID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videoId":   video.ID,
		"title":     video.Title,
		"status":    video.Status,
		"createdAt": video.CreatedAt,
	})
}
