package metricdata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	reqJson, _ := json.MarshalIndent(req, "", "  ")
	log.Println("🚨 RECEIVED RAW REQUEST:")
	log.Println(string(reqJson))

	log.Printf("DEBUG: chart_type=%q, from=%q, to=%q, series_len=%d",
		req.ChartType, req.TimeRange.From, req.TimeRange.To, len(req.Series))

	seriesData, err := s.repo.GetMetricSeriesData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get metric failed: %w", err)
	}

	// Optional forecast/anomaly logic here

	return &metricpb.MetricSeriesResponse{Results: seriesData}, nil
}
