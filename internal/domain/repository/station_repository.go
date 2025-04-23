package repository

import (
	"context"

	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/domain/model"
)

type StationRepository interface {
	FindStationByID(ctx context.Context, id int32) (*model.StationWithLocation, error)
	FilterStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*model.Station, error)
	InsertStation(ctx context.Context, st model.Station) error
	UpdateStationStatus(ctx context.Context, id int32, status string) error

	GetParametersByStationIDs(ctx context.Context, stationIDs []int32) ([]*model.StationParameter, error)
	GetDistinctParametersByStationIDs(ctx context.Context, stationIDs []int32) ([]*model.Parameter, error)

	FindStationIDsByWaterBody(ctx context.Context, waterBodyID int32) ([]int32, error)
	FindStationIDsByCatchment(ctx context.Context, catchmentID int32) ([]int32, error)
	FindStationIDsByRiverBasin(ctx context.Context, riverBasinID int32) ([]int32, error)

	// ========== RIVER BASIN CRUD ==========
	GetRiverBasinByID(ctx context.Context, id int32) (*model.RiverBasin, error)
	ListRiverBasins(ctx context.Context) ([]*model.RiverBasin, error)
	CreateRiverBasin(ctx context.Context, rb model.RiverBasin) error
	DeleteRiverBasin(ctx context.Context, id int32) error
	UpdateRiverBasin(ctx context.Context, rb model.RiverBasin) error
	// ========== WATER BODY CRUD ==========
	GetWaterBodyByID(ctx context.Context, id int32) (*model.WaterBody, error)
	ListWaterBodies(ctx context.Context) ([]*model.WaterBody, error)
	CreateWaterBody(ctx context.Context, wb model.WaterBody) error
	DeleteWaterBody(ctx context.Context, id int32) error
	UpdateWaterBody(ctx context.Context, wb model.WaterBody) error
	// ========== CATCHMENT CRUD ==========
	ListCatchments(ctx context.Context) ([]*model.Catchment, error)
	GetCatchmentByID(ctx context.Context, id int32) (*model.Catchment, error)
	CreateCatchment(ctx context.Context, c model.Catchment) error
	DeleteCatchment(ctx context.Context, id int32) error

	ListEnumValues(table string) ([]*common.EnumValue, error)
}
