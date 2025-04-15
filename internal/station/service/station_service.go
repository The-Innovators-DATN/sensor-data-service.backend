package service

import (
	"context"

	"sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/station/domain"
)

type StationService interface {
	GetStationByID(ctx context.Context, id int32) (*domain.StationWithLocation, error)
	ListStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*domain.Station, error)
	CreateStation(ctx context.Context, st domain.Station) error
	DisableStation(ctx context.Context, id int32) error
	GetStationsByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]int32, error)

	ListEnum(table string) ([]*common.EnumValue, error)
}
