package domain

import "time"

type Catchment struct {
	ID           int32     `json:"id"`
	Name         string    `json:"name"`
	RiverBasinID string    `json:"river_basin_id"` // giữ nguyên là string vì DB để varchar
	Country      string    `json:"country"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
}
