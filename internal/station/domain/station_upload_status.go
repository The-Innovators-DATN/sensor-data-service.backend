package domain

import "time"

type StationUploadStatus struct {
	UID          int32     `json:"uid"`
	AttachmentID int32     `json:"attachment_id"`
	Status       string    `json:"status"`
	UserID       *int32    `json:"user_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
