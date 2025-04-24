package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/domain/model"
	"sensor-data-service.backend/internal/domain/repository"
)

type stationService struct {
	repo repository.StationRepository
}

func NewStationService(repo repository.StationRepository) StationService {
	return &stationService{
		repo: repo,
	}
}

func (s *stationService) GetStationByID(ctx context.Context, id int32) (*model.StationWithLocation, error) {
	st, err := s.repo.FindStationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get station failed: %w", err)
	}
	return st, nil
}

func (s *stationService) ListStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*model.Station, error) {
	return s.repo.FilterStations(ctx, keyword, waterBodyID, catchmentID, riverBasinID)
}

func (s *stationService) CreateStation(ctx context.Context, st model.Station) error {
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
	log.Printf("GetStationsByTarget: targetType=%v, targetID=%d", targetType, targetID)
	switch targetType {
	case stationpb.TargetType_STATION:
		return []int32{targetID}, nil
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

func (s *stationService) GetStationBysByStationType(ctx context.Context, stationType string) ([]*model.Station, error) {
	log.Printf("GetStationBysByStationType: stationType=%s", stationType)
	stations, err := s.repo.GetStationBysByStationType(ctx, stationType)
	if err != nil {
		return nil, fmt.Errorf("get stations by station type failed: %w", err)
	}
	log.Printf("GetStationBysByStationType: stations=%v", stations)
	if len(stations) == 0 {
		return nil, fmt.Errorf("no stations found for station type: %s", stationType)
	}
	return stations, nil
}
func (s *stationService) GetParametersByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]*model.StationParameter, error) {
	log.Printf("GetParametersByTarget: targetType=%v, targetID=%d", targetType, targetID)
	stationIDs, err := s.GetStationsByTarget(ctx, targetType, targetID)
	log.Printf("GetParametersByTarget: stationIDs=%v", stationIDs)
	if err != nil {
		return nil, err
	}
	return s.repo.GetParametersByStationIDs(ctx, stationIDs)
}
func (s *stationService) GetDistinctParametersByTarget(ctx context.Context, targetType stationpb.TargetType, targetID int32) ([]*model.Parameter, error) {
	stationIDs, err := s.GetStationsByTarget(ctx, targetType, targetID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetDistinctParametersByStationIDs(ctx, stationIDs)
}

// RiverBasin CRUD
func (s *stationService) GetRiverBasin(ctx context.Context, id int32) (*model.RiverBasin, error) {
	return s.repo.GetRiverBasinByID(ctx, id)
}
func (s *stationService) GetWaterBody(ctx context.Context, id int32) (*model.WaterBody, error) {
	return s.repo.GetWaterBodyByID(ctx, id)
}

func (s *stationService) ListRiverBasins(ctx context.Context) ([]*model.RiverBasin, error) {
	return s.repo.ListRiverBasins(ctx)
}
func (s *stationService) ListWaterBodies(ctx context.Context) ([]*model.WaterBody, error) {
	return s.repo.ListWaterBodies(ctx)
}
func (s *stationService) CreateRiverBasin(ctx context.Context, rb model.RiverBasin) error {
	if rb.Name == "" {
		return errors.New("missing required fields")
	}
	return s.repo.CreateRiverBasin(ctx, rb)
}
func (s *stationService) CreateWaterBody(ctx context.Context, wb model.WaterBody) error {
	if wb.Name == "" {
		return errors.New("missing required fields")
	}
	return s.repo.CreateWaterBody(ctx, wb)
}
func (s *stationService) DeleteRiverBasin(ctx context.Context, id int32) error {
	return s.repo.DeleteRiverBasin(ctx, id)
}
func (s *stationService) DeleteWaterBody(ctx context.Context, id int32) error {
	return s.repo.DeleteWaterBody(ctx, id)
}
func (s *stationService) UpdateRiverBasin(ctx context.Context, rb model.RiverBasin) error {
	if rb.Name == "" {
		return errors.New("missing required fields")
	}
	return s.repo.UpdateRiverBasin(ctx, rb)
}
func (s *stationService) UpdateWaterBody(ctx context.Context, wb model.WaterBody) error {
	if wb.Name == "" {
		return errors.New("missing required fields")
	}
	return s.repo.UpdateWaterBody(ctx, wb)
}

// Catchments  CRUD
func (s *stationService) ListCatchments(ctx context.Context) ([]*model.Catchment, error) {
	return s.repo.ListCatchments(ctx)
}

func (s *stationService) ListCatchmentByRiverBasin(ctx context.Context, riverBasinID int32) ([]*model.Catchment, error) {
	return s.repo.GetCatchmentsByRiverBasinID(ctx, riverBasinID)
}

func (s *stationService) ListWaterBodyByCatchment(ctx context.Context, catchmentID int32) ([]*model.WaterBody, error) {
	return s.repo.GetWaterBodyByCatchmentID(ctx, catchmentID)
}

func (s *stationService) GetCatchmentByID(ctx context.Context, id int32) (*model.Catchment, error) {
	return s.repo.GetCatchmentByID(ctx, id)
}

func (s *stationService) CreateCatchment(ctx context.Context, c model.Catchment) error {
	return s.repo.CreateCatchment(ctx, c)
}

func (s *stationService) DeleteCatchment(ctx context.Context, id int32) error {
	return s.repo.DeleteCatchment(ctx, id)
}
func (s *stationService) ListEnum(table string) ([]*common.EnumValue, error) {
	return s.repo.ListEnumValues(table)
}
