package metric

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type store struct {
	conn clickhouse.Conn
}

func NewClickhouseStore(conn clickhouse.Conn) Store {
	return &store{conn: conn}
}

func (s *store) ExecQuery(ctx context.Context, sql string, args ...any) ([]map[string]interface{}, error) {
	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	cols := rows.Columns()
	colTypes := rows.ColumnTypes()

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(cols))

		// Gán đúng loại cho từng cột để scan
		for i, ct := range colTypes {
			dbType := ct.DatabaseTypeName()

			switch {
			case dbType == "Float32":
				var v float32
				values[i] = &v
			case dbType == "Float64":
				var v float64
				values[i] = &v
			case dbType == "Int32":
				var v int32
				values[i] = &v
			case dbType == "Int64":
				var v int64
				values[i] = &v
			case dbType == "String":
				var v string
				values[i] = &v
			case dbType == "DateTime" || strings.HasPrefix(dbType, "DateTime64"):
				var v time.Time
				values[i] = &v
			default:
				var v interface{}
				values[i] = &v
			}
		}

		if err := rows.Scan(values...); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			rowMap[col] = deref(values[i])
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func (s *store) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := s.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}
	return nil
}

func deref(ptr interface{}) interface{} {
	switch v := ptr.(type) {
	case *float32:
		return *v
	case *float64:
		return *v
	case *int32:
		return *v
	case *int64:
		return *v
	case *string:
		return *v
	case *time.Time:
		return *v
	default:
		return v
	}
}
