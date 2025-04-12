package metricdata

import (
	"context"
	"fmt"
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

	for _, s := range req.Series {
		if s.TargetType != metricpb.TargetType_STATION {
			continue
		}
		from := strings.TrimSuffix(req.TimeRange.From, "Z")
		to := strings.TrimSuffix(req.TimeRange.To, "Z")

		query := fmt.Sprintf(`
			SELECT
				value,
				metric_id,
				station_id,
				trend_anomaly,
				datetime
			FROM messages_sharded
			WHERE datetime BETWEEN toDateTime('%s') AND toDateTime('%s')
			AND station_id = %d
			AND metric_id = %d
			ORDER BY datetime
		`, from, to, s.TargetId, s.MetricId)

		records, err := r.chdb.ExecQuery(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("query failed for ref_id %s: %w", s.RefId, err)
		}

		var series []*metricpb.MetricPoint
		for _, row := range records {
			dt := row["datetime"].(time.Time)

			// Convert trend_anomaly from *bool
			var trendAnomaly bool
			if b, ok := row["trend_anomaly"].(*bool); ok && b != nil {
				trendAnomaly = *b
			}
			series = append(series, &metricpb.MetricPoint{
				Datetime:     dt.Format(time.RFC3339),
				Value:        float32(toFloat(row["value"])),
				TrendAnomaly: trendAnomaly,
			})
			// log.Printf("series: %v", series)
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

// Safe convert interface{} -> bool
func toBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case uint8:
		return val != 0
	case int:
		return val != 0
	case int64:
		return val != 0
	case float64:
		return val != 0
	case string:
		return val == "1" || val == "true"
	default:
		return false
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
