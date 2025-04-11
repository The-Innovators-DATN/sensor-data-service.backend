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
