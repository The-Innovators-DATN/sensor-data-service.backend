// metricdata_repository.go (refactored)
package repository

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	metricpb "sensor-data-service.backend/api/pb/metricdatapb"
	"sensor-data-service.backend/internal/common/castutil"
	"sensor-data-service.backend/internal/infrastructure/cache"
	"sensor-data-service.backend/internal/infrastructure/db"
	"sensor-data-service.backend/internal/infrastructure/metric"
)

const RedisRetention = 7 * 24 * time.Hour // 7 days

// MetricDataRepository wraps CH + Redis + relational store
type MetricDataRepository struct {
	store db.Store
	cache cache.Store
	chdb  metric.Store
}

func NewMetricDataRepository(ch metric.Store, s db.Store, c cache.Store) *MetricDataRepository {
	return &MetricDataRepository{chdb: ch, store: s, cache: c}
}

// -------------------------------------------------- public entry --------------------------------------------------

func (r *MetricDataRepository) GetMetricSeriesData(ctx context.Context, req *metricpb.MetricSeriesRequest) ([]*metricpb.SeriesData, error) {
	log.Print("[debug] GetMetricSeriesData called")
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if len(req.Series) == 0 {
		return nil, fmt.Errorf("no series requested")
	}
	if req.TimeRange == nil {
		return nil, fmt.Errorf("no time range specified")
	}
	if req.TimeRange.From == "" || req.TimeRange.To == "" {
		return nil, fmt.Errorf("invalid time range")
	}
	step := req.StepSeconds
	if step <= 0 {
		step = 3600 // default 1 h
	}

	from, err := time.Parse(time.RFC3339, req.TimeRange.From)
	if err != nil {
		return nil, fmt.Errorf("invalid time_range.from: %w", err)
	}
	to, err := time.Parse(time.RFC3339, req.TimeRange.To)
	if err != nil {
		return nil, fmt.Errorf("invalid time_range.to: %w", err)
	}

	cutoff := time.Now().UTC().Add(-RedisRetention)
	var out []*metricpb.SeriesData

	for _, sel := range req.Series {
		stations, err := r.resolveStations(ctx, sel.TargetType, sel.TargetId)
		if err != nil {
			return nil, err
		}
		if len(stations) == 0 {
			log.Printf("[warn] no station for target=%v:%d", sel.TargetType, sel.TargetId)
			continue
		}

		var points []*metricpb.MetricPoint
		switch {
		case to.Before(cutoff):
			points, err = r.queryAggregatorCH(ctx, stations, sel.MetricId, from, to, step)
		case from.After(cutoff):
			points, err = r.queryAggregatorRedis(ctx, stations, sel.MetricId, from, to, step)
		default:
			oldPt, err1 := r.queryAggregatorCH(ctx, stations, sel.MetricId, from, cutoff, step)
			if err1 != nil {
				return nil, err1
			}
			newPt, err2 := r.queryAggregatorRedis(ctx, stations, sel.MetricId, cutoff, to, step)
			if err2 != nil {
				return nil, err2
			}
			points = mergeSeriesPoints(oldPt, newPt)
		}
		if err != nil {
			return nil, err
		}

		out = append(out, &metricpb.SeriesData{
			RefId:      sel.RefId,
			TargetType: sel.TargetType,
			TargetId:   sel.TargetId,
			MetricId:   sel.MetricId,
			Series:     points,
		})
	}
	return out, nil
}

// ----------------------------------------------- station resolution ------------------------------------------------

func (r *MetricDataRepository) resolveStations(ctx context.Context, t metricpb.TargetType, id int32) ([]int32, error) {
	if t == metricpb.TargetType_STATION {
		return []int32{id}, nil
	}
	key := fmt.Sprintf("station_targets:%d:%d", t, id)

	var cached []int32
	ok, err := r.cache.GetJSON(ctx, key, &cached)
	if err != nil {
		log.Printf("[error][cache] failed to GetJSON for key=%s: %v", key, err)
	}
	if ok {
		return cached, nil
	}

	sql, err := buildStationQuery(t)
	if err != nil {
		return nil, fmt.Errorf("[resolve][buildQuery] unsupported target type (%v): %w", t, err)
	}

	rows, err := r.store.ExecQuery(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("[resolve][store] failed to exec query (%s): %w", sql, err)
	}

	var ids []int32
	for _, r := range rows {
		ids = append(ids, int32(castutil.ToInt(r["id"])))
	}

	if err := r.cache.SetJSON(ctx, key, ids, int64(time.Hour.Seconds())); err != nil {
		log.Printf("[warn][cache] failed to SetJSON cache for key=%s: %v", key, err)
	}
	return ids, nil
}

func buildStationQuery(t metricpb.TargetType) (string, error) {
	switch t {
	case metricpb.TargetType_WATER_BODY:
		return "SELECT id FROM station WHERE water_body_id = $1", nil
	case metricpb.TargetType_CATCHMENT:
		return `SELECT s.id FROM station s JOIN water_body w ON s.water_body_id = w.id WHERE w.catchment_id = $1`, nil
	case metricpb.TargetType_RIVER_BASIN:
		return `SELECT s.id FROM station s JOIN water_body w ON s.water_body_id = w.id JOIN catchment c ON w.catchment_id = c.id WHERE c.river_basin_id = $1`, nil
	default:
		return "", fmt.Errorf("unsupported target type: %v", t)
	}
}

// ----------------------------------------------- ClickHouse aggregator ---------------------------------------------

func (r *MetricDataRepository) queryAggregatorCH(ctx context.Context, stations []int32, metricID int32, from, to time.Time, step int32) ([]*metricpb.MetricPoint, error) {
	if len(stations) == 0 {
		return nil, nil
	}

	var cond string
	if len(stations) == 1 {
		cond = fmt.Sprintf("station_id = %d", stations[0])
	} else {
		cond = fmt.Sprintf("station_id IN (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(stations)), ","), "[]"))
	}

	q := fmt.Sprintf(`SELECT toStartOfInterval(datetime, INTERVAL %d second) AS t, avg(value) AS v, max(trend_anomaly) AS a FROM messages_sharded WHERE %s AND metric_id = %d AND datetime BETWEEN toDateTime('%s') AND toDateTime('%s') GROUP BY t ORDER BY t`, step, cond, metricID, from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05"))

	rows, err := r.chdb.ExecQuery(ctx, q)
	if err != nil {
		log.Printf("[error][ch-query] failed CH query: %s", q)
		return nil, fmt.Errorf("[ch-query] failed: %w", err)
	}

	var pts []*metricpb.MetricPoint
	for _, r := range rows {
		ts := r["t"].(time.Time)
		val, _ := castutil.ToFloat(r["v"])
		pts = append(pts, &metricpb.MetricPoint{
			Datetime:     ts.Format(time.RFC3339),
			Value:        float32(val),
			TrendAnomaly: castutil.ToBool(r["a"]),
		})
	}
	return pts, nil
}

// ----------------------------------------------- Redis aggregator --------------------------------------------------

func (r *MetricDataRepository) queryAggregatorRedis(ctx context.Context, stations []int32, metricID int32, from, to time.Time, step int32) ([]*metricpb.MetricPoint, error) {
	if len(stations) == 1 {
		return r.singleStationRedis(ctx, stations[0], metricID, from, to, step)
	}
	return r.multiStationRedis(ctx, stations, metricID, from, to, step)
}

func (r *MetricDataRepository) singleStationRedis(ctx context.Context, stationID, metricID int32, from, to time.Time, step int32) ([]*metricpb.MetricPoint, error) {
	tsKey := fmt.Sprintf("sensor_%d_%d", stationID, metricID)
	anomalyKey := fmt.Sprintf("trendanomaly:%d:%d", stationID, metricID)

	ensureSeries(ctx, r.cache, tsKey)
	ensureSet(ctx, r.cache, anomalyKey)

	bucket := time.Duration(step) * time.Second
	data, err := r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucket)
	if isKeyMissing(err) {
		ensureSeries(ctx, r.cache, tsKey)
		err = nil
		data = nil
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		// back-fill nếu Redis chưa có
		if err := r.RefreshRedisSeriesData(ctx, stationID, metricID, from, to); err != nil {
			log.Printf("[warn][redis-refresh] failed to refresh Redis data for station=%d metric=%d: %v", stationID, metricID, err)
		}

		// Retry lần 1
		data, err = r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucket)
		if err != nil && isKeyMissing(err) {
			// Retry lần 2: nếu key vẫn chưa được tạo do CH cũng rỗng → đảm bảo tạo key
			_ = r.cache.TSCreate(ctx, tsKey, RedisRetention)
			data = []cache.TSTimestampValue{}
			err = nil
		}
		if err != nil {
			return nil, fmt.Errorf("TSRangeAgg after refresh failed: %w", err)
		}
	}

	var pts []*metricpb.MetricPoint
	for _, dp := range data {
		ts := time.UnixMilli(int64(dp.Timestamp))
		flag, _ := r.cache.SIsMember(ctx, anomalyKey, fmt.Sprintf("%d", dp.Timestamp))
		pts = append(pts, &metricpb.MetricPoint{
			Datetime:     ts.Format(time.RFC3339),
			Value:        float32(dp.Value),
			TrendAnomaly: flag,
		})
	}
	return pts, nil
}

// --------------------------------------------------------------------------------------------------
// multi‑station redis aggregator (with key‑missing fix) --------------------------------------------

func (r *MetricDataRepository) multiStationRedis(ctx context.Context, stations []int32, metricID int32, from, to time.Time, step int32) ([]*metricpb.MetricPoint, error) {
	type bucket struct {
		sum float64
		n   int
		an  bool
	}
	buckets := map[int64]*bucket{}
	bucketDur := time.Duration(step) * time.Second

	for _, st := range stations {
		tsKey := fmt.Sprintf("sensor_%d_%d", st, metricID)
		anomalyKey := fmt.Sprintf("trendanomaly:%d:%d", st, metricID)
		ensureSeries(ctx, r.cache, tsKey)
		ensureSet(ctx, r.cache, anomalyKey)

		data, err := r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucketDur)
		if isKeyMissing(err) {
			ensureSeries(ctx, r.cache, tsKey)
			err = nil
			data = nil
		}
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			if err := r.RefreshRedisSeriesData(ctx, st, metricID, from, to); err != nil {
				log.Printf("[warn][redis-refresh] failed to refresh Redis data for station=%d metric=%d: %v", st, metricID, err)
			}
			data, _ = r.cache.TSRangeAgg(ctx, tsKey, from, to, cache.Avg, bucketDur)
		}

		for _, dp := range data {
			ts := dp.Timestamp
			b := buckets[ts]
			if b == nil {
				b = &bucket{}
				buckets[ts] = b
			}
			b.sum += dp.Value
			b.n++
			flag, _ := r.cache.SIsMember(ctx, anomalyKey, fmt.Sprintf("%d", ts))
			b.an = b.an || flag
		}
	}

	var keys []int64
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	var pts []*metricpb.MetricPoint
	for _, k := range keys {
		b := buckets[k]
		pts = append(pts, &metricpb.MetricPoint{
			Datetime:     time.UnixMilli(k).Format(time.RFC3339),
			Value:        float32(b.sum / float64(b.n)),
			TrendAnomaly: b.an,
		})
	}
	return pts, nil
}

// ----------------------------------------------- CH → Redis back‑fill ---------------------------------------------

func (r *MetricDataRepository) RefreshRedisSeriesData(ctx context.Context, stationID, metricID int32, from, to time.Time) error {
	ensureSeries(ctx, r.cache, fmt.Sprintf("sensor_%d_%d", stationID, metricID))
	ensureSet(ctx, r.cache, fmt.Sprintf("trendanomaly:%d:%d", stationID, metricID))

	q := fmt.Sprintf(`SELECT toUnixTimestamp(datetime)*1000 AS ts, value, trend_anomaly FROM messages_sharded WHERE station_id=%d AND metric_id=%d AND datetime BETWEEN toDateTime('%s') AND toDateTime('%s')`, stationID, metricID, from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05"))
	rows, err := r.chdb.ExecQuery(ctx, q)
	if err != nil {
		return fmt.Errorf("[refresh][CH] exec query failed for station=%d metric=%d: %w", stationID, metricID, err)
	}

	tsKey := fmt.Sprintf("sensor_%d_%d", stationID, metricID)
	anKey := fmt.Sprintf("trendanomaly:%d:%d", stationID, metricID)
	anomalySetKey := fmt.Sprintf("trendanomaly:%d:%d", stationID, metricID)
	if len(rows) == 0 {
		_ = r.cache.TSCreate(ctx, tsKey, RedisRetention)
		_ = r.cache.Set(ctx, anomalySetKey, "", 0) // Set rỗng để không bị lỗi SIsMember
	}
	for _, rRow := range rows {
		tsF, _ := castutil.ToFloat(rRow["ts"])
		ts := int64(tsF)
		val, _ := castutil.ToFloat(rRow["value"])
		_ = r.cache.TSAdd(ctx, tsKey, time.UnixMilli(ts), val)
		if castutil.ToBool(rRow["trend_anomaly"]) {
			_ = r.cache.SAdd(ctx, anKey, fmt.Sprintf("%d", ts))
		} else {
			_ = r.cache.SRem(ctx, anKey, fmt.Sprintf("%d", ts))
		}
	}
	return nil
}

// ----------------------------------------------- helpers -----------------------------------------------------------
// --------------------------------------------------------------------------------------------------
// util helpers -------------------------------------------------------------------------------------

func isKeyMissing(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "key does not exist") || strings.Contains(err.Error(), "TSDB")
}

func ensureSeries(ctx context.Context, c cache.Store, key string) {
	_ = c.TSCreate(ctx, key, RedisRetention) // create if absent, ignore error otherwise
}

func ensureSet(ctx context.Context, c cache.Store, key string) {
	_ = c.Set(ctx, key, "", 0) // create empty string key if not exist
}

func mergeSeriesPoints(a, b []*metricpb.MetricPoint) []*metricpb.MetricPoint {
	merged := append(a, b...)
	sort.Slice(merged, func(i, j int) bool { return merged[i].Datetime < merged[j].Datetime })
	return merged
}
