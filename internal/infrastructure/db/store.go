package db

import (
	"context"
)

type Store interface {
	ExecQuery(ctx context.Context, sql string, args ...any) ([]map[string]interface{}, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
}
