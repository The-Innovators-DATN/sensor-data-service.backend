package transport

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
	"sensor-data-service.backend/api/pb/commonpb"
	"sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/station/domain"
	"sensor-data-service.backend/internal/station/service"
)

type StationHandler struct {
	stationpb.UnimplementedStationServiceServer
	svc service.StationService
}

func NewStationHandler(svc service.StationService) *StationHandler {
	return &StationHandler{
		svc: svc,
	}
}

// ===== Station =====

func (h *StationHandler) GetStation(ctx context.Context, req *stationpb.StationID) (*commonpb.StandardResponse, error) {
	result, err := h.svc.GetStationByID(ctx, req.Id)
	if err != nil {
		return common.WrapError(fmt.Sprintf("Failed to get station: %v", err)), nil
	}

	station := &stationpb.Station{
		Id:             result.Station.ID,
		Name:           result.Station.Name,
		Description:    result.Station.Description,
		Lat:            result.Station.Lat,
		Long:           result.Station.Long,
		Status:         result.Station.Status,
		StationType:    result.Station.StationType,
		Country:        result.Station.Country,
		WaterBodyId:    result.Station.WaterBodyID,
		StationManager: result.Station.StationManager,
		CreatedAt:      result.Station.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      result.Station.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	location := convertToProtoLocation(result.Location)

	return common.WrapSuccess("GetStation OK", &stationpb.StationWithLocation{
		Station:  station,
		Location: location,
	})
}

func convertToProtoLocation(loc domain.StationLocation) *stationpb.StationLocation {
	return &stationpb.StationLocation{
		WaterBodyName:  loc.WaterBodyName,
		WaterBodyType:  loc.WaterBodyType,
		CatchmentId:    loc.CatchmentID,
		CatchmentName:  loc.CatchmentName,
		CatchmentDesc:  loc.CatchmentDesc,
		RiverBasinId:   loc.RiverBasinID,
		RiverBasinName: loc.RiverBasinName,
	}
}

func (h *StationHandler) ListStations(ctx context.Context, req *stationpb.StationQuery) (*commonpb.StandardResponse, error) {
	stations, err := h.svc.ListStations(ctx, req.Keyword, req.WaterBodyId, req.CatchmentId, req.RiverBasinId)
	if err != nil {
		return common.WrapError(fmt.Sprintf("ListStations failed: %v", err)), nil
	}

	var pbStations []*stationpb.Station
	for _, s := range stations {
		pbStations = append(pbStations, &stationpb.Station{
			Id:             s.ID,
			Name:           s.Name,
			Description:    s.Description,
			Lat:            s.Lat,
			Long:           s.Long,
			Status:         s.Status,
			StationType:    s.StationType,
			Country:        s.Country,
			WaterBodyId:    s.WaterBodyID,
			StationManager: s.StationManager,
			CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return common.WrapSuccess("ListStations OK", &stationpb.StationList{Stations: pbStations})
}

func (h *StationHandler) CreateStation(ctx context.Context, req *stationpb.Station) (*commonpb.StandardResponse, error) {
	entity := domain.Station{
		Name:           req.Name,
		Description:    req.Description,
		Lat:            req.Lat,
		Long:           req.Long,
		Status:         req.Status,
		StationType:    req.StationType,
		Country:        req.Country,
		WaterBodyID:    req.WaterBodyId,
		StationManager: req.StationManager,
	}

	err := h.svc.CreateStation(ctx, entity)
	if err != nil {
		return common.WrapError(fmt.Sprintf("CreateStation failed: %v", err)), nil
	}

	return common.WrapSuccess("Station created", &stationpb.Station{
		Name:           entity.Name,
		Description:    entity.Description,
		Lat:            entity.Lat,
		Long:           entity.Long,
		Status:         entity.Status,
		StationType:    entity.StationType,
		Country:        entity.Country,
		WaterBodyId:    entity.WaterBodyID,
		StationManager: entity.StationManager,
	})
}

func (h *StationHandler) DisableStation(ctx context.Context, req *stationpb.StationID) (*commonpb.StandardResponse, error) {
	err := h.svc.DisableStation(ctx, req.Id)
	if err != nil {
		return common.WrapError(fmt.Sprintf("DisableStation failed: %v", err)), nil
	}
	return common.WrapSuccess("Station disabled", &emptypb.Empty{})
}

func (h *StationHandler) GetStationsByTarget(ctx context.Context, req *stationpb.TargetSelector) (*commonpb.StandardResponse, error) {
	stations, err := h.svc.GetStationsByTarget(ctx, req.TargetType, req.TargetId)
	if err != nil {
		return common.WrapError(fmt.Sprintf("GetStationsByTarget failed: %v", err)), nil
	}
	return common.WrapSuccess("GetStationsByTarget OK", &stationpb.StationIDList{
		StationIds: stations,
	})
}

// ===== ENUMS =====

func (h *StationHandler) ListCountries(ctx context.Context, _ *emptypb.Empty) (*commonpb.StandardResponse, error) {
	values, err := h.svc.ListEnum("country")
	if err != nil {
		return common.WrapError("ListCountries failed"), nil
	}
	var pbVals []*stationpb.EnumValue
	for _, v := range values {
		pbVals = append(pbVals, &stationpb.EnumValue{Name: v.Name})
	}
	return common.WrapSuccess("ListCountries OK", &stationpb.EnumList{Values: pbVals})
}

func (h *StationHandler) ListStationTypes(ctx context.Context, _ *emptypb.Empty) (*commonpb.StandardResponse, error) {
	values, err := h.svc.ListEnum("station_type")
	if err != nil {
		return common.WrapError("ListStationTypes failed"), nil
	}
	var pbVals []*stationpb.EnumValue
	for _, v := range values {
		pbVals = append(pbVals, &stationpb.EnumValue{Name: v.Name})
	}
	return common.WrapSuccess("ListStationTypes OK", &stationpb.EnumList{Values: pbVals})
}

func (h *StationHandler) ListStatus(ctx context.Context, _ *emptypb.Empty) (*commonpb.StandardResponse, error) {
	values, err := h.svc.ListEnum("status")
	if err != nil {
		return common.WrapError("ListStatus failed"), nil
	}
	var pbVals []*stationpb.EnumValue
	for _, v := range values {
		pbVals = append(pbVals, &stationpb.EnumValue{Name: v.Name})
	}
	return common.WrapSuccess("ListStatus OK", &stationpb.EnumList{Values: pbVals})
}
