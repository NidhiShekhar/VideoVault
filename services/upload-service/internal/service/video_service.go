package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/nidhi/video-vault/services/upload-service/internal/model"
)

type UploadVideoInput struct {
	Title       string `validate:"required,min=3,max=120"`
	Description string `validate:"max=500"`
	Filename    string `validate:"required"`
}

type VideoRepository interface {
	Save(ctx context.Context, video model.Video) error
	GetByID(ctx context.Context, videoID string) (model.Video, error)
}

type VideoService struct {
	repo      VideoRepository
	validator *validator.Validate
}

func NewVideoService(repo VideoRepository) *VideoService {
	return &VideoService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *VideoService) Upload(ctx context.Context, in UploadVideoInput) (model.Video, error) {
	if err := s.validator.Struct(in); err != nil {
		return model.Video{}, fmt.Errorf("validate upload input: %w", err)
	}

	video := model.Video{
		ID:          uuid.NewString(),
		Title:       in.Title,
		Description: in.Description,
		Filename:    in.Filename,
		Status:      model.StatusPending,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Save(ctx, video); err != nil {
		return model.Video{}, fmt.Errorf("save video: %w", err)
	}

	return video, nil
}

func (s *VideoService) GetStatus(ctx context.Context, videoID string) (model.Video, error) {
	video, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return model.Video{}, fmt.Errorf("get video by id: %w", err)
	}

	return video, nil
}
