package station

type LocationType string

const (
	Station    LocationType = "station"
	RiverBasin LocationType = "river_basin"
	Catchment  LocationType = "catchment"
	WaterBody  LocationType = "water_body"
	Country    LocationType = "country"
)

type Location struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
