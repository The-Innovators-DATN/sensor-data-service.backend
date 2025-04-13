package cache

import (
	"context"
	"encoding/json"
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
func (r *redisStore) TSCreate(ctx context.Context, key string, retention time.Duration) error {
	opts := &redis.TSOptions{
		Retention: retention.Milliseconds(),
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

func (r *redisStore) TSRangeAgg(
	ctx context.Context,
	key string,
	from, to time.Time,
	agg redis.Aggregator, // like "avg", "min", "max"
	bucketDuration time.Duration, // size of each aggregation bucket
) ([]redis.TSTimestampValue, error) {
	opts := &redis.TSRangeOptions{
		Aggregation:    agg,
		BucketDuration: bucketDuration.Milliseconds(),
	}

	res := r.client.TSRangeWithArgs(ctx, key, int(from.UnixMilli()), int(to.UnixMilli()), opts)
	return res.Val(), res.Err()
}
func (r *redisStore) SIsMember(ctx context.Context, setKey, member string) (bool, error) {
	return r.client.SIsMember(ctx, setKey, member).Bool()
}
func (r *redisStore) SAdd(ctx context.Context, key string, members ...string) error {
	return r.client.SAdd(ctx, key, members).Err()
}
