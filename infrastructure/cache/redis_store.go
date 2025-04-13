package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) Store {
	return &redisStore{client: client}
}

func (r *redisStore) Set(ctx context.Context, key string, value interface{}, ttlSeconds int64) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *redisStore) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (r *redisStore) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisStore) SetJSON(ctx context.Context, key string, value interface{}, expiration int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, string(data), expiration)
}

func (r *redisStore) GetJSON(ctx context.Context, key string, out interface{}) (bool, error) {
	raw, err := r.Get(ctx, key)
	if err != nil || raw == nil {
		return false, err
	}
	strVal, ok := raw.(string)
	if !ok {
		return false, nil
	}
	err = json.Unmarshal([]byte(strVal), out)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *redisStore) SRem(ctx context.Context, key string, members ...string) error {
	return r.client.SRem(ctx, key, members).Err()
}

func (r *redisStore) TSCreate(ctx context.Context, key string, retention time.Duration) error {
	opts := &redis.TSOptions{
		Retention: int(retention.Milliseconds()), // fix: cast to int
	}
	return r.client.TSCreateWithArgs(ctx, key, opts).Err()
}

func (r *redisStore) TSAdd(ctx context.Context, key string, timestamp time.Time, value float64) error {
	return r.client.TSAdd(ctx, key, timestamp.UnixMilli(), value).Err()
}

func (r *redisStore) TSRange(ctx context.Context, key string, from, to time.Time) ([]redis.TSTimestampValue, error) {
	res := r.client.TSRange(ctx, key, int(from.UnixMilli()), int(to.UnixMilli()))
	return res.Val(), res.Err()
}

// redis.Avg
// redis.Min
// redis.Max
// redis.Sum
// redis.Count
// redis.First
// redis.Last
// redis.StdP
// redis.StdS
// redis.VarP
// redis.VarS

type Aggregator string

const (
	Avg Aggregator = "avg"
	Min Aggregator = "min"
	Max Aggregator = "max"
)

// TSTimestampValue là struct mà TSRangeWithArgs trả về
type TSTimestampValue struct {
	Timestamp int64
	Value     float64
}

type TSRangeOptions struct {
	AggregationType string
	BucketSizeSec   int64
	// ... tuỳ lib
}

func (r *redisStore) TSRangeAgg(
	ctx context.Context,
	key string,
	from, to time.Time,
	agg Aggregator,
	bucketDuration time.Duration,
) ([]TSTimestampValue, error) {
	cmd := r.client.Do(ctx,
		"TS.RANGE",
		key,
		from.UnixMilli(),
		to.UnixMilli(),
		"AGGREGATION",
		string(agg),
		int64(bucketDuration.Milliseconds()),
	)

	vals, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	rawData, ok := vals.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected TS.RANGE response format")
	}

	var out []TSTimestampValue
	for _, item := range rawData {
		point, ok := item.([]interface{})
		if !ok || len(point) != 2 {
			continue
		}

		tsInt, ok1 := point[0].(int64)
		valFloat, ok2 := point[1].(float64)
		if !ok1 {
			if tsStr, ok := point[0].(string); ok {
				parsed, err := strconv.ParseInt(tsStr, 10, 64)
				if err != nil {
					fmt.Printf("[warn] cannot parse timestamp: %v\n", err)
					continue
				}
				tsInt = parsed
				ok1 = true
			}
		}
		if ok1 && ok2 {
			out = append(out, TSTimestampValue{
				Timestamp: tsInt,
				Value:     valFloat,
			})
		}
	}
	return out, nil
}

// SIsMember kiểm tra member trong set => trả về bool
func (r *redisStore) SIsMember(ctx context.Context, setKey, member string) (bool, error) {
	val, err := r.client.SIsMember(ctx, setKey, member).Result()
	if err != nil {
		return false, err
	}
	return val, nil
}
func (r *redisStore) SAdd(ctx context.Context, key string, members ...string) error {
	return r.client.SAdd(ctx, key, members).Err()
}
