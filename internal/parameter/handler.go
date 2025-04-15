package parameter

import (
	"context"
	"time"

	commonpb "sensor-data-service.backend/api/pb/commonpb"
	pb "sensor-data-service.backend/api/pb/parameterpb"
	"sensor-data-service.backend/internal/common"
)

type Handler struct {
	service *Service
	pb.UnimplementedParameterServiceServer
}

func NewGrpcHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListParameters(ctx context.Context, _ *pb.Empty) (*commonpb.StandardResponse, error) {
	params, err := h.service.ListParameters(ctx)
	if err != nil {
		return common.WrapError("failed to fetch parameters"), nil
	}

	var res []*pb.ParameterResponse
	for _, p := range params {
		res = append(res, convertToProto(p))
	}
	return common.WrapSuccess("retrieved successfully", &pb.ParameterListResponse{Parameters: res})
}

func (h *Handler) GetParameter(ctx context.Context, req *pb.ParameterRequest) (*commonpb.StandardResponse, error) {
	param, err := h.service.GetParameter(ctx, int(req.Id))
	if err != nil {
		return common.WrapError("parameter not found"), nil
	}
	return common.WrapSuccess("retrieved successfully", convertToProto(param))
}

func (h *Handler) CreateParameter(ctx context.Context, req *pb.ParameterCreateRequest) (*commonpb.StandardResponse, error) {
	now := time.Now()
	param := Parameter{
		Name:           req.Name,
		Unit:           req.Unit,
		ParameterGroup: req.ParameterGroup,
		Description:    req.Description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := h.service.CreateParameter(ctx, param); err != nil {
		return common.WrapError("failed to create parameter"), nil
	}
	return common.WrapSuccess("created successfully", convertToProto(param))
}

func (h *Handler) UpdateParameter(ctx context.Context, req *pb.ParameterUpdateRequest) (*commonpb.StandardResponse, error) {
	param := Parameter{
		ID:             int(req.Id),
		Name:           req.Name,
		Unit:           req.Unit,
		ParameterGroup: req.ParameterGroup,
		Description:    req.Description,
		UpdatedAt:      time.Now(),
	}
	if err := h.service.UpdateParameter(ctx, param); err != nil {
		return common.WrapError("failed to update parameter"), nil
	}
	return common.WrapSuccess("updated successfully", convertToProto(param))
}

func (h *Handler) DeleteParameter(ctx context.Context, req *pb.ParameterRequest) (*commonpb.StandardResponse, error) {
	if err := h.service.DeleteParameter(ctx, int(req.Id)); err != nil {
		return common.WrapError("failed to delete parameter"), nil
	}
	// Nếu bạn muốn trả message đơn giản thì define một proto struct DeleteMessage:
	// message DeleteMessage { string status = 1; }
	return common.WrapSuccess("deleted successfully", &pb.DeleteResponse{Status: "deleted"})
}

// Helper to convert internal model to proto response
func convertToProto(p Parameter) *pb.ParameterResponse {
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
