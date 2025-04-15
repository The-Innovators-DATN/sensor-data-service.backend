package domain

import "time"

type StationParameter struct {
	StationID     int32      `json:"station_id"`
	ParameterID   int32      `json:"parameter_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastReceiveAt *time.Time `json:"last_receive_at,omitempty"`
	LastValue     *float64   `json:"last_value,omitempty"`
	Status        string     `json:"status"` // enum: active, inactive, deleted
}
