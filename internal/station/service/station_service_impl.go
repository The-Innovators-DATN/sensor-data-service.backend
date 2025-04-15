package service

import (
	"context"
	"errors"
	"fmt"

	"sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/station/domain"
	"sensor-data-service.backend/internal/station/repository"
)

type stationService struct {
	repo repository.StationRepository
}

func NewStationService(repo repository.StationRepository) StationService {
	return &stationService{
		repo: repo,
	}
}

func (s *stationService) GetStationByID(ctx context.Context, id int32) (*domain.StationWithLocation, error) {
	st, err := s.repo.FindStationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get station failed: %w", err)
	}
	return st, nil
}

func (s *stationService) ListStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*domain.Station, error) {
	return s.repo.FilterStations(ctx, keyword, waterBodyID, catchmentID, riverBasinID)
}

func (s *stationService) CreateStation(ctx context.Context, st domain.Station) error {
	// simple validation
	if st.Name == "" || st.Status == "" || st.Country == "" {
		return errors.New("missing required fields")
	}
	return s.repo.InsertStation(ctx, st)
}

func (s *stationService) DisableStation(ctx context.Context, id int32) error {
	return s.repo.UpdateStationStatus(ctx, id, "inactive")
}

func (s *stationService) GetStationsByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]int32, error) {
	switch targetType {
	case stationpb.TargetType_WATER_BODY:
		return s.repo.FindStationIDsByWaterBody(ctx, targetID)
	case stationpb.TargetType_CATCHMENT:
		return s.repo.FindStationIDsByCatchment(ctx, targetID)
	case stationpb.TargetType_RIVER_BASIN:
		return s.repo.FindStationIDsByRiverBasin(ctx, targetID)
	default:
		return nil, errors.New("unsupported target type")
	}
}

func (s *stationService) ListEnum(table string) ([]*common.EnumValue, error) {
	return s.repo.ListEnumValues(table)
}
