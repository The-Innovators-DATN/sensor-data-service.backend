package metricdata

import (
	"context"
	"fmt"

	metricpb "sensor-data-service.backend/api/pb/metricdatapb"
)

type MetricDataService struct {
	repo *MetricDataRepository
}

func NewMetricDataService(repo *MetricDataRepository) *MetricDataService {
	return &MetricDataService{repo: repo}
}

// GetMetricSeries processes the MetricSeriesRequest and returns the results.
func (s *MetricDataService) GetMetricSeries(ctx context.Context, req *metricpb.MetricSeriesRequest) (*metricpb.MetricSeriesResponse, error) {

	seriesData, err := s.repo.GetMetricSeriesData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get metric failed: %w", err)
	}

	// Optional forecast/anomaly logic here

	return &metricpb.MetricSeriesResponse{Results: seriesData}, nil
}
