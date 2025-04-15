package domain

import (
	"time"

	"github.com/google/uuid"
)

type StarDashboard struct {
	ID                  uuid.UUID `json:"id"`
	UserID              int32     `json:"user_id"`
	Version             int       `json:"version"`
	LayoutConfiguration string    `json:"layout_configuration"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
