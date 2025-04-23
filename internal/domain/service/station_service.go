package service

import (
	"context"

	"sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/domain/model"
)

type StationService interface {
	GetStationByID(ctx context.Context, id int32) (*model.StationWithLocation, error)
	ListStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*model.Station, error)
	CreateStation(ctx context.Context, st model.Station) error
	DisableStation(ctx context.Context, id int32) error

	GetStationsByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]int32, error)
	GetParametersByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]*model.StationParameter, error)
	GetDistinctParametersByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]*model.Parameter, error)

	// GetStationAttachments(ctx context.Context, stationID int32) ([]*model.StationAttachment, error)
	GetRiverBasin(ctx context.Context, id int32) (*model.RiverBasin, error)
	GetWaterBody(ctx context.Context, id int32) (*model.WaterBody, error)
	ListRiverBasins(ctx context.Context) ([]*model.RiverBasin, error)
	ListWaterBodies(ctx context.Context) ([]*model.WaterBody, error)
	// ListCatchmentByRiverBasin(ctx context.Context, riverBasinID int32) ([]*model.Catchment, error)
	// ListWaterBodyByCatchment(ctx context.Context, catchmentID int32) ([]*model.WaterBody, error)
	// ListWaterBodyByRiverBasin(ctx context.Context, riverBasinID int32) ([]*model.WaterBody, error)
	CreateRiverBasin(ctx context.Context, rb model.RiverBasin) error
	CreateWaterBody(ctx context.Context, wb model.WaterBody) error
	DeleteRiverBasin(ctx context.Context, id int32) error
	DeleteWaterBody(ctx context.Context, id int32) error
	UpdateRiverBasin(ctx context.Context, rb model.RiverBasin) error
	UpdateWaterBody(ctx context.Context, wb model.WaterBody) error

	ListCatchments(ctx context.Context) ([]*model.Catchment, error)
	CreateCatchment(ctx context.Context, c model.Catchment) error
	DeleteCatchment(ctx context.Context, id int32) error

	// GetRiverBasin(ctx context.Context, id int32) (*model.RiverBasin, error)
	GetCatchmentByID(ctx context.Context, id int32) (*model.Catchment, error)
	// GetWaterBody(ctx context.Context, id int32) (*model.WaterBody, error)
	ListEnum(table string) ([]*common.EnumValue, error)
}
