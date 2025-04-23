package handler

import (
	"context"
	"fmt"
	"log"

	commonpb "sensor-data-service.backend/api/pb/commonpb"
	metricdatapb "sensor-data-service.backend/api/pb/metricdatapb"
	common "sensor-data-service.backend/internal/common"
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

	data, err := h.service.GetMetricSeries(ctx, req)
	if err != nil {
		return common.WrapError(fmt.Sprintf("get metric failed: %v", err)), nil
	}

	return common.WrapSuccess("retrieved successfully", data)
}
