package station

type LocationType string

const (
	Station      LocationType = "station"
	RiverBasin   LocationType = "river_basin"
	Catchment    LocationType = "catchment"
	WaterBody    LocationType = "water_body"
	Country      LocationType = "country"
)
type Location struct {