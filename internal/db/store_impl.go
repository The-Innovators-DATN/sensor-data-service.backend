package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type store struct {
	conn *pgx.Conn
}

func NewPostgresStore(conn *pgx.Conn) Store {
	return &store{conn: conn}
}

func (s *store) ExecQuery(ctx context.Context, sql string, args ...any) ([]map[string]interface{}, error) {
	rows, err := s.conn.Query(ctx, sql, args...)
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

	return results, nil
}
