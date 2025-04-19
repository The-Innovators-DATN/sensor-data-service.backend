package service

import (
	"context"

	"sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/parameter"
	"sensor-data-service.backend/internal/station/domain"
)

type StationService interface {
	GetStationByID(ctx context.Context, id int32) (*domain.StationWithLocation, error)
	ListStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*domain.Station, error)
	CreateStation(ctx context.Context, st domain.Station) error
	DisableStation(ctx context.Context, id int32) error

	GetStationsByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]int32, error)
	GetParametersByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]*domain.StationParameter, error)
	GetDistinctParametersByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]*parameter.Parameter, error)

	// GetStationAttachments(ctx context.Context, stationID int32) ([]*domain.StationAttachment, error)
	GetRiverBasin(ctx context.Context, id int32) (*domain.RiverBasin, error)
	GetWaterBody(ctx context.Context, id int32) (*domain.WaterBody, error)
	ListRiverBasins(ctx context.Context) ([]*domain.RiverBasin, error)
	ListWaterBodies(ctx context.Context) ([]*domain.WaterBody, error)
	// ListCatchmentByRiverBasin(ctx context.Context, riverBasinID int32) ([]*domain.Catchment, error)
	// ListWaterBodyByCatchment(ctx context.Context, catchmentID int32) ([]*domain.WaterBody, error)
	// ListWaterBodyByRiverBasin(ctx context.Context, riverBasinID int32) ([]*domain.WaterBody, error)
	CreateRiverBasin(ctx context.Context, rb domain.RiverBasin) error
	CreateWaterBody(ctx context.Context, wb domain.WaterBody) error
	DeleteRiverBasin(ctx context.Context, id int32) error
	DeleteWaterBody(ctx context.Context, id int32) error
	UpdateRiverBasin(ctx context.Context, rb domain.RiverBasin) error
	UpdateWaterBody(ctx context.Context, wb domain.WaterBody) error

	ListCatchments(ctx context.Context) ([]*domain.Catchment, error)
	CreateCatchment(ctx context.Context, c domain.Catchment) error
	DeleteCatchment(ctx context.Context, id int32) error

	// GetRiverBasin(ctx context.Context, id int32) (*domain.RiverBasin, error)
	GetCatchmentByID(ctx context.Context, id int32) (*domain.Catchment, error)
	// GetWaterBody(ctx context.Context, id int32) (*domain.WaterBody, error)
	ListEnum(table string) ([]*common.EnumValue, error)
}
