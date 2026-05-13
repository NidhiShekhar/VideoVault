package model

import "time"

type UploadStatus string

const (
	StatusPending UploadStatus = "pending"
)

type Video struct {
	ID          string
	Title       string
	Description string
	Filename    string
	Status      UploadStatus
	CreatedAt   time.Time
}
