package model

import (
	"time"

	"github.com/google/uuid"
)

type Dashboard struct {
	UID                 uuid.UUID `json:"uid"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	LayoutConfiguration string    `json:"layout_configuration"`
	CreatedBy           int32     `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Version             int32     `json:"version"`
	Status              string    `json:"status"`
}

type PaginatedDashboards struct {
	Items []*Dashboard
	Total int32
}
