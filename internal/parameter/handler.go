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
		res = append(res, &pb.ParameterResponse{
			Id:             int32(p.ID),
			Name:           p.Name,
			Unit:           p.Unit,
			ParameterGroup: p.ParameterGroup,
			Description:    p.Description,
			CreatedAt:      p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ParameterListResponse{Parameters: res}, nil
}
