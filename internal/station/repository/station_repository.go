package repository

import (
	"context"

	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/station/domain"
)

type StationRepository interface {
	FindStationByID(ctx context.Context, id int32) (*domain.StationWithLocation, error)
	FilterStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*domain.Station, error)
	InsertStation(ctx context.Context, st domain.Station) error
	UpdateStationStatus(ctx context.Context, id int32, status string) error

	FindStationIDsByWaterBody(ctx context.Context, waterBodyID int32) ([]int32, error)
	FindStationIDsByCatchment(ctx context.Context, catchmentID int32) ([]int32, error)
	FindStationIDsByRiverBasin(ctx context.Context, riverBasinID int32) ([]int32, error)

	ListEnumValues(table string) ([]*common.EnumValue, error)
}
