package metricdata

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	metricpb "sensor-data-service.backend/api/pb/metricdatapb"
	"sensor-data-service.backend/infrastructure/cache"
	"sensor-data-service.backend/infrastructure/db"
	"sensor-data-service.backend/infrastructure/metric"
)

// Thời hạn mà data nằm trong RedisTimeSeries
const RedisRetention = 7 * 24 * time.Hour // 7 ngày

type MetricDataRepository struct {
	store db.Store
	cache cache.Store
	chdb  metric.Store
}

func NewMetricDataRepository(chdb metric.Store, store db.Store, cache cache.Store) *MetricDataRepository {
	return &MetricDataRepository{chdb: chdb, store: store, cache: cache}
}

// GetMetricSeriesData: hàm chính xử lý request
func (r *MetricDataRepository) GetMetricSeriesData(ctx context.Context, req *metricpb.MetricSeriesRequest) ([]*metricpb.SeriesData, error) {
	var results []*metricpb.SeriesData

	stepSeconds := req.StepSeconds
	if stepSeconds <= 0 {
		stepSeconds = 3600 // fallback 1h
	}

	from, err := time.Parse(time.RFC3339, req.TimeRange.From)
	if err != nil {
		return nil, fmt.Errorf("invalid time_range.from: %w", err)
	}
	to, err := time.Parse(time.RFC3339, req.TimeRange.To)
	if err != nil {
		return nil, fmt.Errorf("invalid time_range.to: %w", err)
	}

	now := time.Now().UTC()
	cutoff := now.Add(-RedisRetention)

	for _, sel := range req.Series {
		// Lấy danh sách station(s) tuỳ target_type
		var stationIDs []int32
		if sel.TargetType == metricpb.TargetType_STATION {
			stationIDs = []int32{sel.TargetId}
		} else {
			list, err := r.store.GetStationsByTarget(ctx, sel.TargetType, sel.TargetId)
			if err != nil {
				return nil, fmt.Errorf("getStationsByTarget failed: %w", err)
			}
			if len(list) == 0 {
				log.Printf("[warn] no stations for target_type=%v, target_id=%v", sel.TargetType, sel.TargetId)
				continue
			}
			stationIDs = list
		}

		var seriesPoints []*metricpb.MetricPoint
		switch {
		case to.Before(cutoff):
			// Hoàn toàn ngoài 7 ngày => Query CH
			seriesPoints, err = r.queryAggregatorCH(ctx, stationIDs, sel.MetricId, from, to, stepSeconds)
			if err != nil {
				return nil, err
			}
		case from.After(cutoff):
			// Hoàn toàn trong 7 ngày => Query Redis
			seriesPoints, err = r.queryAggregatorRedis(ctx, stationIDs, sel.MetricId, from, to, stepSeconds)
			if err != nil {
				return nil, err
			}
		default:
			// Giao nhau => tách 2 phần
			chTo := cutoff
			redisFrom := cutoff

			oldPoints, err1 := r.queryAggregatorCH(ctx, stationIDs, sel.MetricId, from, chTo, stepSeconds)
			if err1 != nil {
				return nil, err1
			}
			newPoints, err2 := r.queryAggregatorRedis(ctx, stationIDs, sel.MetricId, redisFrom, to, stepSeconds)
			if err2 != nil {
				return nil, err2
			}
			seriesPoints = mergeSeriesPoints(oldPoints, newPoints)
		}

		results = append(results, &metricpb.SeriesData{
			RefId:      sel.RefId,
			TargetType: sel.TargetType,
			TargetId:   sel.TargetId,
			MetricId:   sel.MetricId,
			Series:     seriesPoints,
		})
	}

	return results, nil
}

// ------------------- Aggregator CH -------------------

// queryAggregatorCH: group by station(s) + interval => avg(value), maxOrNull(trend_anomaly)
func (r *MetricDataRepository) queryAggregatorCH(
	ctx context.Context,
	stationIDs []int32,
	metricID int32,
	from, to time.Time,
	stepSeconds int32,
) ([]*metricpb.MetricPoint, error) {

	if len(stationIDs) == 0 {
		return nil, nil
	}

	var query string
	if len(stationIDs) == 1 {
		query = fmt.Sprintf(`
SELECT
    toStartOfInterval(datetime, INTERVAL %d second) AS interval_time,
    avg(value) AS avg_val,
    maxOrNull(trend_anomaly) AS anomaly
FROM messages_sharded
WHERE station_id = %d
  AND metric_id = %d
  AND datetime BETWEEN toDateTime('%s') AND toDateTime('%s')
GROUP BY interval_time
ORDER BY interval_time`,
			stepSeconds, stationIDs[0], metricID,
			from.Format("2006-01-02 15:04:05"),
			to.Format("2006-01-02 15:04:05"),
		)
	} else {
		idStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(stationIDs)), ","), "[]")
		query = fmt.Sprintf(`
SELECT
    toStartOfInterval(datetime, INTERVAL %d second) AS interval_time,
    avg(value) AS avg_val,
    maxOrNull(trend_anomaly) AS anomaly
FROM messages_sharded
WHERE station_id IN (%s)
  AND metric_id = %d
  AND datetime BETWEEN toDateTime('%s') AND toDateTime('%s')
GROUP BY interval_time
ORDER BY interval_time`,
			stepSeconds, idStr, metricID,
			from.Format("2006-01-02 15:04:05"),
			to.Format("2006-01-02 15:04:05"),
		)
	}

	rows, err := r.chdb.ExecQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("clickhouse aggregator query failed: %w", err)
	}

	var points []*metricpb.MetricPoint
	for _, row := range rows {
		itime, _ := row["interval_time"].(time.Time)
		avgVal, _ := toFloat(row["avg_val"])
		anomaly := toBool(row["anomaly"])

		points = append(points, &metricpb.MetricPoint{
			Datetime:     itime.Format(time.RFC3339),
			Value:        float32(avgVal),
			TrendAnomaly: anomaly,
		})
	}
	return points, nil
}

// ------------------- Aggregator Redis -------------------

// queryAggregatorRedis: nếu 1 station => queryRedisAggregatorSingle
// nếu nhiều station => aggregateAcrossStationsRedis
func (r *MetricDataRepository) queryAggregatorRedis(
	ctx context.Context,
	stationIDs []int32,
	metricID int32,
	from, to time.Time,
	stepSeconds int32,
) ([]*metricpb.MetricPoint, error) {

	if len(stationIDs) == 0 {
		return nil, nil
	}
	if len(stationIDs) == 1 {
		return r.queryRedisAggregatorSingle(ctx, stationIDs[0], metricID, from, to, stepSeconds)
	}
	return r.aggregateAcrossStationsRedis(ctx, stationIDs, metricID, from, to, stepSeconds)
}

// queryRedisAggregatorSingle: aggregator 1 station => TSRangeAgg.
// Nếu trả về rỗng => gọi RefreshRedisSeriesData() => aggregator lần nữa.
func (r *MetricDataRepository) queryRedisAggregatorSingle(
	ctx context.Context,
	stationID int32,
	metricID int32,
	from, to time.Time,
	stepSeconds int32,
) ([]*metricpb.MetricPoint, error) {

	tsKey := fmt.Sprintf("sensor_%d_%d", stationID, metricID)
	anomalySetKey := fmt.Sprintf("trendanomaly:%d:%d", stationID, metricID)

	bucketDur := time.Duration(stepSeconds) * time.Second

	aggData, err := r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucketDur)
	if err != nil {
		return nil, fmt.Errorf("TSRangeAgg fail key=%s: %w", tsKey, err)
	}

	// Nếu chưa có data => refresh Redis từ CH
	if len(aggData) == 0 {
		if err := r.RefreshRedisSeriesData(ctx, stationID, metricID, from, to); err != nil {
			log.Printf("[warn] RefreshRedisSeriesData failed st=%d,metric=%d => %v", stationID, metricID, err)
		}
		// Gọi lại aggregator
		aggData, err = r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucketDur)
		if err != nil {
			return nil, fmt.Errorf("TSRangeAgg after refresh fail key=%s: %w", tsKey, err)
		}
	}

	var points []*metricpb.MetricPoint
	for _, dp := range aggData {
		tsMillis := dp.Timestamp
		tParsed := time.UnixMilli(int64(tsMillis)).UTC()

		exist, err := r.cache.SIsMember(ctx, anomalySetKey, fmt.Sprintf("%d", tsMillis))
		if err != nil {
			// fallback
			exist = false
			log.Printf("[warn] SIsMember fail key=%s => %v", anomalySetKey, err)
		}

		points = append(points, &metricpb.MetricPoint{
			Datetime:     tParsed.Format(time.RFC3339),
			Value:        float32(dp.Value),
			TrendAnomaly: exist,
		})
	}
	return points, nil
}

// aggregateAcrossStationsRedis: aggregator cho nhiều station => loop từng station => merge
func (r *MetricDataRepository) aggregateAcrossStationsRedis(
	ctx context.Context,
	stationIDs []int32,
	metricID int32,
	from, to time.Time,
	stepSeconds int32,
) ([]*metricpb.MetricPoint, error) {

	type bucketAgg struct {
		sum     float64
		count   int
		anomaly bool
	}
	aggMap := make(map[int64]*bucketAgg)
	bucketDur := time.Duration(stepSeconds) * time.Second

	for _, stID := range stationIDs {
		tsKey := fmt.Sprintf("sensor_%d_%d", stID, metricID)
		anomalySetKey := fmt.Sprintf("trendanomaly:%d:%d", stID, metricID)

		// Lấy data
		stationAgg, err := r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucketDur)
		if err != nil {
			return nil, fmt.Errorf("TSRangeAgg fail key=%s: %w", tsKey, err)
		}

		// Nếu rỗng => refresh -> query lại
		if len(stationAgg) == 0 {
			if err := r.RefreshRedisSeriesData(ctx, stID, metricID, from, to); err != nil {
				log.Printf("[warn] RefreshRedisSeriesData fail st=%d,m=%d => %v", stID, metricID, err)
			}
			stationAgg, err = r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucketDur)
			if err != nil {
				return nil, fmt.Errorf("TSRangeAgg after refresh fail key=%s: %w", tsKey, err)
			}
		}

		// Gom data vào aggMap
		for _, dp := range stationAgg {
			tsMillis := dp.Timestamp

			exist, err := r.cache.SIsMember(ctx, anomalySetKey, fmt.Sprintf("%d", tsMillis))
			if err != nil {
				exist = false
				log.Printf("[warn] SIsMember fail key=%s => %v", anomalySetKey, err)
			}

			if b, ok := aggMap[tsMillis]; ok {
				b.sum += dp.Value
				b.count++
				b.anomaly = b.anomaly || exist
			} else {
				aggMap[tsMillis] = &bucketAgg{
					sum:     dp.Value,
					count:   1,
					anomaly: exist,
				}
			}
		}
	}

	// map -> sorted slice
	var timestamps []int64
	for ts := range aggMap {
		timestamps = append(timestamps, ts)
	}
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] < timestamps[j]
	})

	var points []*metricpb.MetricPoint
	for _, ts := range timestamps {
		entry := aggMap[ts]
		avgVal := entry.sum / float64(entry.count)
		tParsed := time.UnixMilli(ts).UTC()

		points = append(points, &metricpb.MetricPoint{
			Datetime:     tParsed.Format(time.RFC3339),
			Value:        float32(avgVal),
			TrendAnomaly: entry.anomaly,
		})
	}

	return points, nil
}

// ------------------- Refresh Redis (CH -> Redis) -------------------

// RefreshRedisSeriesData: Lấy data [from,to] từ CH, ghi vào Redis TS + anomaly set
// - TS.ADD sensor_{stationID}_{metricID} {timestamp_ms} {value}
// - if anomaly => SADD trendanomaly:stationID:metricID {timestamp_ms} else => SREM ...
func (r *MetricDataRepository) RefreshRedisSeriesData(
	ctx context.Context,
	stationID, metricID int32,
	from, to time.Time,
) error {

	// Query CH data + anomaly
	query := fmt.Sprintf(`
SELECT
    toUnixTimestamp(datetime)*1000 AS ts_ms,
    value,
    trend_anomaly
FROM messages_sharded
WHERE station_id = %d
  AND metric_id = %d
  AND datetime BETWEEN toDateTime('%s') AND toDateTime('%s')
ORDER BY datetime
`,
		stationID,
		metricID,
		from.Format("2006-01-02 15:04:05"),
		to.Format("2006-01-02 15:04:05"),
	)

	rows, err := r.chdb.ExecQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("RefreshRedisSeriesData: CH query failed: %w", err)
	}

	tsKey := fmt.Sprintf("sensor_%d_%d", stationID, metricID)
	anomalySetKey := fmt.Sprintf("trendanomaly:%d:%d", stationID, metricID)

	for _, row := range rows {
		tsMs := int64(toFloat(row["ts_ms"]))
		val, _ := toFloat(row["value"])
		anomaly := toBool(row["trend_anomaly"])

		// Ghi data point vào Redis TimeSeries
		// TS.ADD <key> <timestamp_ms> <value>
		if err := r.cache.TSAdd(ctx, tsKey, time.UnixMilli(tsMs), val); err != nil {
			log.Printf("[warn] TSAdd fail key=%s, ts=%d => %v", tsKey, tsMs, err)
		}

		// Update anomaly set
		tsStr := fmt.Sprintf("%d", tsMs)
		if anomaly {
			if err := r.cache.SAdd(ctx, anomalySetKey, tsStr); err != nil {
				log.Printf("[warn] SAdd fail for key=%s => %v", anomalySetKey, err)
			}
		}

	}
	return nil
}

// ------------------- Utils -------------------

func mergeSeriesPoints(oldList, newList []*metricpb.MetricPoint) []*metricpb.MetricPoint {
	merged := append(oldList, newList...)
	sort.Slice(merged, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, merged[i].Datetime)
		tj, _ := time.Parse(time.RFC3339, merged[j].Datetime)
		return ti.Before(tj)
	})
	return merged
}

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

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
