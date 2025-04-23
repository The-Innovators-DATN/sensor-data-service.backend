package model

import "time"

type StationKey struct {
	ID        int32     `json:"id"`
	StationID int32     `json:"station_id"`
	OrgID     int32     `json:"org_id"`
	IsRevoked bool      `json:"is_revoked"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
