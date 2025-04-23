package model

import (
	"time"

	"github.com/google/uuid"
)

type StationAttachment struct {
	UID           uuid.UUID `json:"uid"`
	Size          int       `json:"size"`
	Filename      string    `json:"filename"`
	ContentType   string    `json:"content_type"`
	DisplayName   string    `json:"display_name"`
	WorkflowState string    `json:"workflow_state"`
	FileState     string    `json:"file_state"`
	Namespace     string    `json:"namespace"`
	StationID     int32     `json:"station_id"`
	UserID        int32     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
