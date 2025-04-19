package repository

import (
	"context"

	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/parameter"
	"sensor-data-service.backend/internal/station/domain"
)

type StationRepository interface {
	FindStationByID(ctx context.Context, id int32) (*domain.StationWithLocation, error)
	FilterStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*domain.Station, error)
	InsertStation(ctx context.Context, st domain.Station) error
	UpdateStationStatus(ctx context.Context, id int32, status string) error

	GetParametersByStationIDs(ctx context.Context, stationIDs []int32) ([]*domain.StationParameter, error)
	GetDistinctParametersByStationIDs(ctx context.Context, stationIDs []int32) ([]*parameter.Parameter, error)

	FindStationIDsByWaterBody(ctx context.Context, waterBodyID int32) ([]int32, error)
	FindStationIDsByCatchment(ctx context.Context, catchmentID int32) ([]int32, error)
	FindStationIDsByRiverBasin(ctx context.Context, riverBasinID int32) ([]int32, error)

	// ========== RIVER BASIN CRUD ==========
	GetRiverBasinByID(ctx context.Context, id int32) (*domain.RiverBasin, error)
	ListRiverBasins(ctx context.Context) ([]*domain.RiverBasin, error)
	CreateRiverBasin(ctx context.Context, rb domain.RiverBasin) error
	DeleteRiverBasin(ctx context.Context, id int32) error
	UpdateRiverBasin(ctx context.Context, rb domain.RiverBasin) error
	// ========== WATER BODY CRUD ==========
	GetWaterBodyByID(ctx context.Context, id int32) (*domain.WaterBody, error)
	ListWaterBodies(ctx context.Context) ([]*domain.WaterBody, error)
	CreateWaterBody(ctx context.Context, wb domain.WaterBody) error
	DeleteWaterBody(ctx context.Context, id int32) error
	UpdateWaterBody(ctx context.Context, wb domain.WaterBody) error
	// ========== CATCHMENT CRUD ==========
	ListCatchments(ctx context.Context) ([]*domain.Catchment, error)
	GetCatchmentByID(ctx context.Context, id int32) (*domain.Catchment, error)
	CreateCatchment(ctx context.Context, c domain.Catchment) error
	DeleteCatchment(ctx context.Context, id int32) error

	ListEnumValues(table string) ([]*common.EnumValue, error)
}
