package transport

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
	"sensor-data-service.backend/api/pb/commonpb"
	"sensor-data-service.backend/api/pb/parameterpb"
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

func (h *StationHandler) ListRiverBasins(ctx context.Context, _ *emptypb.Empty) (*commonpb.StandardResponse, error) {
	data, err := h.svc.ListRiverBasins(ctx)
	if err != nil {
		return common.WrapError("failed to list river basins"), nil
	}
	var res []*stationpb.RiverBasin
	for _, r := range data {
		res = append(res, &stationpb.RiverBasin{
			Id:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Status:      r.Status,
			UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return common.WrapSuccess("ok", &stationpb.RiverBasinList{RiverBasins: res})
}
func (h *StationHandler) GetRiverBasinByID(ctx context.Context, req *stationpb.RiverBasinID) (*commonpb.StandardResponse, error) {
	r, err := h.svc.GetRiverBasin(ctx, req.Id)
	if err != nil {
		return common.WrapError("not found"), nil
	}
	return common.WrapSuccess("ok", &stationpb.RiverBasin{
		Id:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
func (h *StationHandler) CreateRiverBasin(ctx context.Context, req *stationpb.RiverBasin) (*commonpb.StandardResponse, error) {
	r := domain.RiverBasin{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	err := h.svc.CreateRiverBasin(ctx, r)
	if err != nil {
		return common.WrapError(fmt.Sprintf("create failed: %v", err)), nil
	}
	return common.WrapSuccess("created", &emptypb.Empty{})
}

func (h *StationHandler) UpdateRiverBasin(ctx context.Context, req *stationpb.RiverBasin) (*commonpb.StandardResponse, error) {
	r := domain.RiverBasin{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	err := h.svc.UpdateRiverBasin(ctx, r)
	if err != nil {
		return common.WrapError(fmt.Sprintf("update failed: %v", err)), nil
	}
	return common.WrapSuccess("updated", &emptypb.Empty{})
}

func (h *StationHandler) DeleteRiverBasin(ctx context.Context, req *stationpb.RiverBasinID) (*commonpb.StandardResponse, error) {
	err := h.svc.DeleteRiverBasin(ctx, req.Id)
	if err != nil {
		return common.WrapError("delete failed"), nil
	}
	return common.WrapSuccess("deleted", &emptypb.Empty{})
}

func (h *StationHandler) ListWaterBodies(ctx context.Context, _ *emptypb.Empty) (*commonpb.StandardResponse, error) {
	data, err := h.svc.ListWaterBodies(ctx)
	if err != nil {
		return common.WrapError("failed to list water bodies"), nil
	}
	var res []*stationpb.WaterBody
	for _, w := range data {
		res = append(res, &stationpb.WaterBody{
			Id:          w.ID,
			Name:        w.Name,
			Description: w.Description,
			Status:      w.Status,
			UpdatedAt:   w.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return common.WrapSuccess("ok", &stationpb.WaterBodyList{WaterBodies: res})
}

func (h *StationHandler) GetWaterBodyByID(ctx context.Context, req *stationpb.WaterBodyID) (*commonpb.StandardResponse, error) {
	w, err := h.svc.GetWaterBody(ctx, req.Id)
	if err != nil {
		return common.WrapError("not found"), nil
	}
	return common.WrapSuccess("ok", &stationpb.WaterBody{
		Id:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		Status:      w.Status,
		UpdatedAt:   w.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *StationHandler) CreateWaterBody(ctx context.Context, req *stationpb.WaterBody) (*commonpb.StandardResponse, error) {
	w := domain.WaterBody{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	err := h.svc.CreateWaterBody(ctx, w)
	if err != nil {
		return common.WrapError(fmt.Sprintf("create failed: %v", err)), nil
	}
	return common.WrapSuccess("created", &emptypb.Empty{})
}

func (h *StationHandler) UpdateWaterBody(ctx context.Context, req *stationpb.WaterBody) (*commonpb.StandardResponse, error) {
	w := domain.WaterBody{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	err := h.svc.UpdateWaterBody(ctx, w)
	if err != nil {
		return common.WrapError(fmt.Sprintf("update failed: %v", err)), nil
	}
	return common.WrapSuccess("updated", &emptypb.Empty{})
}

func (h *StationHandler) DeleteWaterBody(ctx context.Context, req *stationpb.WaterBodyID) (*commonpb.StandardResponse, error) {
	err := h.svc.DeleteWaterBody(ctx, req.Id)
	if err != nil {
		return common.WrapError("delete failed"), nil
	}
	return common.WrapSuccess("deleted", &emptypb.Empty{})
}

func (h *StationHandler) ListCatchments(ctx context.Context, req *stationpb.CatchmentQuery) (*commonpb.StandardResponse, error) {
	data, err := h.svc.ListCatchments(ctx)
	if err != nil {
		return common.WrapError("failed to list catchments"), nil
	}
	var res []*stationpb.Catchment
	for _, c := range data {
		res = append(res, &stationpb.Catchment{
			Id:           c.ID,
			Name:         c.Name,
			Description:  c.Description,
			Status:       c.Status,
			RiverBasinId: c.RiverBasinID,
			Country:      c.Country,
			UpdatedAt:    c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return common.WrapSuccess("ok", &stationpb.CatchmentList{Catchments: res})
}

func (h *StationHandler) GetCatchmentByID(ctx context.Context, req *stationpb.CatchmentID) (*commonpb.StandardResponse, error) {
	c, err := h.svc.GetCatchmentByID(ctx, req.Id)
	if err != nil {
		return common.WrapError("not found"), nil
	}
	return common.WrapSuccess("ok", &stationpb.Catchment{
		Id:           c.ID,
		Name:         c.Name,
		Description:  c.Description,
		Status:       c.Status,
		RiverBasinId: c.RiverBasinID,
		Country:      c.Country,
		UpdatedAt:    c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *StationHandler) CreateCatchment(ctx context.Context, req *stationpb.Catchment) (*commonpb.StandardResponse, error) {
	c := domain.Catchment{
		Name:         req.Name,
		Description:  req.Description,
		Status:       req.Status,
		RiverBasinID: req.RiverBasinId,
		Country:      req.Country,
	}
	err := h.svc.CreateCatchment(ctx, c)
	if err != nil {
		return common.WrapError(fmt.Sprintf("create failed: %v", err)), nil
	}
	return common.WrapSuccess("created", &emptypb.Empty{})
}
func (h *StationHandler) GetParametersByTarget(ctx context.Context, req *stationpb.TargetSelector) (*commonpb.StandardResponse, error) {
	log.Printf("GetParametersByTarget: targetType=%v, targetID=%d", req.TargetType, req.TargetId)

	params, err := h.svc.GetDistinctParametersByTarget(ctx, req.TargetType, req.TargetId)
	if err != nil {
		return common.WrapError(fmt.Sprintf("failed to get parameters: %v", err)), nil
	}

	var res []*parameterpb.ParameterResponse
	for _, p := range params {
		res = append(res, &parameterpb.ParameterResponse{
			Id:             int32(p.ID),
			Name:           p.Name,
			Unit:           p.Unit,
			ParameterGroup: p.ParameterGroup,
			Description:    p.Description,
			Status:         p.Status,
			CreatedAt:      p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return common.WrapSuccess("ok", &parameterpb.ParameterListResponse{Parameters: res})
}

// func (h *StationHandler) GetParametersByTarget(ctx context.Context, req *stationpb.TargetSelector) (*commonpb.StandardResponse, error) {
// 	log.Printf("GetParametersByTarget: targetType=%v, targetID=%d", req.TargetType, req.TargetId)
// 	params, err := h.svc.GetParametersByTarget(ctx, req.TargetType, req.TargetId)
// 	if err != nil {
// 		return common.WrapError(fmt.Sprintf("failed to get parameters: %v", err)), nil
// 	}

// 	var res []*stationpb.StationParameter
// 	for _, p := range params {
// 		param := &stationpb.StationParameter{
// 			StationId:   p.StationID,
// 			ParameterId: p.ParameterID,
// 			Status:      p.Status,
// 			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
// 		}

// 		if p.LastValue != nil {
// 			param.LastValue = *p.LastValue
// 		}
// 		if p.LastReceiveAt != nil {
// 			param.LastReceivAt = p.LastReceiveAt.Format("2006-01-02T15:04:05Z")
// 		}

// 		res = append(res, param)
// 	}

// 	return common.WrapSuccess("ok", &stationpb.StationParameterList{Items: res})
// }

func (h *StationHandler) DeleteCatchment(ctx context.Context, req *stationpb.CatchmentID) (*commonpb.StandardResponse, error) {
	err := h.svc.DeleteCatchment(ctx, req.Id)
	if err != nil {
		return common.WrapError("delete failed"), nil
	}
	return common.WrapSuccess("deleted", &emptypb.Empty{})
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
