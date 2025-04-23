package metric

import (
	"context"
)

type Store interface {
	ExecQuery(ctx context.Context, sql string, args ...any) ([]map[string]interface{}, error)
	Exec(ctx context.Context, sql string, args ...any) error
}
