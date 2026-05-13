package service_test

import (
	"context"
	"testing"

	"github.com/nidhi/video-vault/services/upload-service/internal/repository/memory"
	"github.com/nidhi/video-vault/services/upload-service/internal/service"
)

func TestVideoServiceUpload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   service.UploadVideoInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: service.UploadVideoInput{
				Title:       "Go Kafka Deep Dive",
				Description: "week 1 test upload",
				Filename:    "video.mp4",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			input: service.UploadVideoInput{
				Title:       "",
				Description: "invalid case",
				Filename:    "video.mp4",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := memory.NewVideoRepository()
			svc := service.NewVideoService(repo)

			video, err := svc.Upload(context.Background(), tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && video.ID == "" {
				t.Fatalf("expected non-empty video ID")
			}
		})
	}
}
