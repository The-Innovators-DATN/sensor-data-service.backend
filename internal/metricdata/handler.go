// ✅ FIXED HANDLER - handler.go
package metricdata

import (
	"context"

	metricdatapb "sensor-data-service.backend/api/pb/metricdatapb"
)

// MetricDataHandler implements gRPC server
type MetricDataHandler struct {
	metricdatapb.UnimplementedMetricDataServiceServer
	service *MetricDataService
}

func NewMetricDataHandler(service *MetricDataService) *MetricDataHandler {
	return &MetricDataHandler{service: service}
}

func (h *MetricDataHandler) GetMetricSeries(ctx context.Context, req *metricdatapb.MetricSeriesRequest) (*metricdatapb.MetricSeriesResponse, error) {
	return h.service.GetMetricSeries(ctx, req)
}
