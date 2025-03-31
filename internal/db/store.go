package db

import (
	"context"
)

type Store interface {
	ExecQuery(ctx context.Context, sql string, args ...any) ([]map[string]interface{}, error)
}
