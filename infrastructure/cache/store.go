package cache

import (
	"context"
)

type Store interface {
	// Set sets a value in the cache with the given key and expiration time.
	Set(ctx context.Context, key string, value interface{}, expiration int64) error
	// Get retrieves a value from the cache by its key.
	Get(ctx context.Context, key string) (interface{}, error)
	// Delete removes a value from the cache by its key.
	Delete(ctx context.Context, key string) error
	SetJSON(ctx context.Context, key string, value interface{}, expiration int64) error
	GetJSON(ctx context.Context, key string, out interface{}) (bool, error)
}
