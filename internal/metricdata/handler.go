// ✅ FIXED HANDLER - handler.go
package metricdata

import (
	"context"
	"fmt"

	metricdatapb "sensor-data-service.backend/api/pb/metricdatapb"
	commonpb "sensor-data-service.backend/api/pb/commonpb"
	common "sensor-data-service.backend/internal/common"
)

// MetricDataHandler implements gRPC server
type MetricDataHandler struct {
	metricdatapb.UnimplementedMetricDataServiceServer
	service *MetricDataService
}

func NewMetricDataHandler(service *MetricDataService) *MetricDataHandler {
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
