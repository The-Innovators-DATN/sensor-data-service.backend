package parameter

import (
	"context"
	"time"

	"sensor-data-service.backend/infrastructure/db"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Parameter, error)
	// ... other CRUD
}

type postgresRepo struct {
	store db.Store
}

func NewRepository(store db.Store) Repository {
	return &postgresRepo{store: store}
}

func (r *postgresRepo) GetAll(ctx context.Context) ([]Parameter, error) {
	sql := `SELECT id, name, unit, parameter_group, description, created_at, updated_at FROM parameter`
	rows, err := r.store.ExecQuery(ctx, sql)
	if err != nil {
		return nil, err
	}

	var result []Parameter
	for _, row := range rows {
		p := Parameter{
			ID:             int(row["id"].(int32)),
			Name:           row["name"].(string),
			Unit:           row["unit"].(string),
			ParameterGroup: row["parameter_group"].(string),
			Description:    row["description"].(string),
			CreatedAt:      row["created_at"].(time.Time),
			UpdatedAt:      row["updated_at"].(time.Time),
		}
		result = append(result, p)
	}
	return result, nil
}
