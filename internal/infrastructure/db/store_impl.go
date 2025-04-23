package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) Store {
	return &store{pool: pool}
}

func (s *store) ExecQuery(ctx context.Context, sql string, args ...any) ([]map[string]interface{}, error) {
	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = string(fd.Name)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("row values error: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (s *store) Exec(ctx context.Context, query string, args ...any) error {
	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}
	return nil
}
