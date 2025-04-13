package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Set(ctx context.Context, key string, value interface{}, expiration int64) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
	SetJSON(ctx context.Context, key string, value interface{}, expiration int64) error
	GetJSON(ctx context.Context, key string, out interface{}) (bool, error)

	TSCreate(ctx context.Context, key string, retention time.Duration) error
	TSAdd(ctx context.Context, key string, timestamp time.Time, value float64) error
	TSRange(ctx context.Context, key string, from, to time.Time) ([]redis.TSTimestampValue, error)

	// Check if the key exists in the Redis set
	SAdd(ctx context.Context, key string, members ...string) error

	SIsMember(ctx context.Context, setKey, member string) (bool, error)
}
