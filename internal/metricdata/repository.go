package metricdata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	metricpb "sensor-data-service.backend/api/pb/metricdatapb"
	"sensor-data-service.backend/infrastructure/cache"
	"sensor-data-service.backend/infrastructure/db"
	"sensor-data-service.backend/infrastructure/metric"
)

type MetricDataRepository struct {
	store db.Store
	cache cache.Store
	chdb  metric.Store
}

func NewMetricDataRepository(chdb metric.Store, store db.Store, cache cache.Store) *MetricDataRepository {
	return &MetricDataRepository{chdb: chdb, store: store, cache: cache}
}
func (r *MetricDataRepository) GetMetricSeriesData(ctx context.Context, req *metricpb.MetricSeriesRequest) ([]*metricpb.SeriesData, error) {
	var results []*metricpb.SeriesData
	reqJson, _ := json.MarshalIndent(req, "", "  ")
	log.Println("🚨 RECEIVED RAW REQUEST:")
	log.Println(string(reqJson))

	log.Printf("DEBUG: chart_type=%q, from=%q, to=%q, series_len=%d",
		req.ChartType, req.TimeRange.From, req.TimeRange.To, len(req.Series))

	for _, s := range req.Series {
		if s.TargetType != metricpb.TargetType_STATION {
			continue
		}
		from := strings.TrimSuffix(req.TimeRange.From, "Z")
		to := strings.TrimSuffix(req.TimeRange.To, "Z")

		query := fmt.Sprintf(`
			SELECT
				value,
				metric,
				station_id,
				local_error,
				datetime
			FROM sensors_to_kafka
			WHERE datetime BETWEEN toDateTime('%s') AND toDateTime('%s')
			AND station_id = %d
			AND metric = %d
			ORDER BY datetime
		`, from, to, s.TargetId, s.MetricId)

		records, err := r.chdb.ExecQuery(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("query failed for ref_id %s: %w", s.RefId, err)
		}

		var series []*metricpb.MetricPoint
		for _, row := range records {
			dt := row["datetime"].(time.Time)
			series = append(series, &metricpb.MetricPoint{
				Datetime:   dt.Format(time.RFC3339),
				Value:      float32(toFloat(row["value"])),
				LocalError: float32(toFloat(row["local_error"])),
			})
		}

		results = append(results, &metricpb.SeriesData{
			RefId:      s.RefId,
			TargetType: s.TargetType,
			TargetId:   s.TargetId,
			MetricId:   s.MetricId,
			Series:     series,
			Forecast:   []*metricpb.MetricPoint{},
		})
	}

	return results, nil
}

// Helper chuyển kiểu an toàn
func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}

func toFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float32:
		return float64(x)
	case float64:
		return x
	default:
		return 0.0
	}
}

// Tìm refId cho cặp station-metric (nếu có)
func findRefID(series []SeriesSelector, stationID, metricID int32) string {
	for _, s := range series {
		if s.TargetID == stationID && s.MetricID == metricID {
			return s.RefID
		}
	}
	return ""
}
