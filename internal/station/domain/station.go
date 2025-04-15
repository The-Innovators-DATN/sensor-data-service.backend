package domain

import "time"

type Station struct {
	ID             int32     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Lat            float32   `json:"lat"`
	Long           float32   `json:"long"`
	Status         string    `json:"status"`
	StationType    string    `json:"station_type"`
	Country        string    `json:"country"`
	WaterBodyID    int32     `json:"water_body_id"`
	StationManager int32     `json:"station_manager"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StationLocation struct {
	WaterBodyName  string
	WaterBodyType  string
	CatchmentID    int32
	CatchmentName  string
	CatchmentDesc  string
	RiverBasinID   int32
	RiverBasinName string
}

type StationWithLocation struct {
	Station  Station
	Location StationLocation
}
