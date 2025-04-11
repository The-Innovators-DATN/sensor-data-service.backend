package parameter

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "sensor-data-service.backend/api/pb/parameterpb"
)

type Handler struct {
	service *Service
	pb.UnimplementedParameterServiceServer
}

func NewGrpcHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListParameters(ctx context.Context, _ *pb.Empty) (*pb.ParameterListResponse, error) {
	params, err := h.service.ListParameters(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch parameters: %v", err)
	}

	var res []*pb.ParameterResponse
	for _, p := range params {
		res = append(res, convertToProto(p))
	}
	return &pb.ParameterListResponse{Parameters: res}, nil
}

func (h *Handler) GetParameter(ctx context.Context, req *pb.ParameterRequest) (*pb.ParameterResponse, error) {
	param, err := h.service.GetParameter(ctx, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "parameter not found: %v", err)
	}
	return convertToProto(param), nil
}

func (h *Handler) CreateParameter(ctx context.Context, req *pb.ParameterCreateRequest) (*pb.ParameterResponse, error) {
	now := time.Now()
	param := Parameter{
		Name:           req.Name,
		Unit:           req.Unit,
		ParameterGroup: req.ParameterGroup,
		Description:    req.Description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	err := h.service.CreateParameter(ctx, param)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create parameter: %v", err)
	}
	return convertToProto(param), nil
}

func (h *Handler) UpdateParameter(ctx context.Context, req *pb.ParameterUpdateRequest) (*pb.ParameterResponse, error) {
	param := Parameter{
		ID:             int(req.Id),
		Name:           req.Name,
		Unit:           req.Unit,
		ParameterGroup: req.ParameterGroup,
		Description:    req.Description,
		UpdatedAt:      time.Now(),
	}
	err := h.service.UpdateParameter(ctx, param)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update parameter: %v", err)
	}
	return convertToProto(param), nil
}

func (h *Handler) DeleteParameter(ctx context.Context, req *pb.ParameterRequest) (*pb.DeleteResponse, error) {
	err := h.service.DeleteParameter(ctx, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete parameter: %v", err)
	}
	return &pb.DeleteResponse{Status: "deleted"}, nil
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
