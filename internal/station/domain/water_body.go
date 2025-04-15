package domain

import "time"

type WaterBody struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	CatchmentID int32     `json:"catchment_id"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
