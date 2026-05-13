package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/nidhi/video-vault/services/upload-service/internal/model"
)

type VideoRepository struct {
	mu     sync.RWMutex
	videos map[string]model.Video
}

func NewVideoRepository() *VideoRepository {
	return &VideoRepository{
		videos: make(map[string]model.Video),
	}
}

func (r *VideoRepository) Save(_ context.Context, video model.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.videos[video.ID] = video
	return nil
}

func (r *VideoRepository) GetByID(_ context.Context, videoID string) (model.Video, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	video, ok := r.videos[videoID]
	if !ok {
		return model.Video{}, fmt.Errorf("video %q not found", videoID)
	}

	return video, nil
}
