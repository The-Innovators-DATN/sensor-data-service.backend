package handler

import (
	"context"
	"time"

	commonpb "sensor-data-service.backend/api/pb/commonpb"
	pb "sensor-data-service.backend/api/pb/parameterpb"
	"sensor-data-service.backend/internal/common/response"
	"sensor-data-service.backend/internal/domain/model"
	"sensor-data-service.backend/internal/domain/service"
)

type ParameterHandler struct {
	service *service.ParameterService
	pb.UnimplementedParameterServiceServer
}

func NewGrpcParameterHandler(service *service.ParameterService) *ParameterHandler {
	return &ParameterHandler{service: service}
}

func (h *ParameterHandler) ListParameters(ctx context.Context, _ *pb.Empty) (*commonpb.StandardResponse, error) {
	params, err := h.service.ListParameters(ctx)
	if err != nil {
		return response.WrapError("failed to fetch parameters", nil)
	}
	var res []*pb.ParameterResponse
	for _, p := range params {
		res = append(res, convertToProto(p))
	}
	return response.WrapSuccess("retrieved successfully", &pb.ParameterListResponse{Parameters: res})
}

func (h *ParameterHandler) GetParameter(ctx context.Context, req *pb.ParameterRequest) (*commonpb.StandardResponse, error) {
	param, err := h.service.GetParameter(ctx, int(req.Id))
	if err != nil {
		return response.WrapError("parameter not found", nil)
	}
	return response.WrapSuccess("retrieved successfully", convertToProto(param))
}

func (h *ParameterHandler) CreateParameter(ctx context.Context, req *pb.ParameterCreateRequest) (*commonpb.StandardResponse, error) {
	now := time.Now()
	param := model.Parameter{
		Name:           req.Name,
		Unit:           req.Unit,
		ParameterGroup: req.ParameterGroup,
		Description:    req.Description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := h.service.CreateParameter(ctx, param); err != nil {
		return response.WrapError("failed to create parameter", nil)
	}
	return response.WrapSuccess("created successfully", convertToProto(param))
}

func (h *ParameterHandler) UpdateParameter(ctx context.Context, req *pb.ParameterUpdateRequest) (*commonpb.StandardResponse, error) {
	param := model.Parameter{
		ID:             int(req.Id),
		Name:           req.Name,
		Unit:           req.Unit,
		ParameterGroup: req.ParameterGroup,
		Description:    req.Description,
		UpdatedAt:      time.Now(),
	}
	if err := h.service.UpdateParameter(ctx, param); err != nil {
		return response.WrapError("failed to update parameter", nil)
	}
	return response.WrapSuccess("updated successfully", convertToProto(param))
}

func (h *ParameterHandler) DeleteParameter(ctx context.Context, req *pb.ParameterRequest) (*commonpb.StandardResponse, error) {
	if err := h.service.DeleteParameter(ctx, int(req.Id)); err != nil {
		return response.WrapError("failed to delete parameter", nil)
	}
	return response.WrapSuccess("deleted successfully", &pb.DeleteResponse{Status: "deleted"})
}

// Helper to convert internal model to proto response
func convertToProto(p model.Parameter) *pb.ParameterResponse {
	return &pb.ParameterResponse{
		Id:             int32(p.ID),
		Name:           p.Name,
		Unit:           p.Unit,
		ParameterGroup: p.ParameterGroup,
		Description:    p.Description,
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
}
