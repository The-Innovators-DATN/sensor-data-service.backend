package handler

import (
	"context"
	"fmt"
	"log"

	commonpb "sensor-data-service.backend/api/pb/commonpb"
	metricdatapb "sensor-data-service.backend/api/pb/metricdatapb"
	"sensor-data-service.backend/internal/common/response"
	"sensor-data-service.backend/internal/domain/service"
)

// MetricDataHandler implements gRPC server
type MetricDataHandler struct {
	metricdatapb.UnimplementedMetricDataServiceServer
	service *service.MetricDataService
}

func NewMetricDataHandler(service *service.MetricDataService) *MetricDataHandler {
	log.Printf("MetricDataHandler initialized")
	return &MetricDataHandler{service: service}
}

func (h *MetricDataHandler) GetMetricSeries(
	ctx context.Context,
	req *metricdatapb.MetricSeriesRequest,
) (*commonpb.StandardResponse, error) {
	// log raw request
	log.Printf("GetMetricSeries request: %v", req)

	data, err := h.service.GetMetricSeries(ctx, req)
	if err != nil {
		return response.WrapError(fmt.Sprintf("get metric failed: %v", err), nil)
	}
	return response.WrapSuccess("retrieved successfully", data)
}
