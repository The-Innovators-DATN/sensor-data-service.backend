package parameter

import (
	"context"
	"fmt"
	"time"

	"sensor-data-service.backend/infrastructure/cache"
	"sensor-data-service.backend/infrastructure/db"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Parameter, error)
	GetByID(ctx context.Context, id int) (Parameter, error)
	Create(ctx context.Context, p Parameter) error
	Update(ctx context.Context, p Parameter) error
	Delete(ctx context.Context, id int) error
}

type postgresRepo struct {
	store db.Store
	cache cache.Store
}

func NewRepository(store db.Store, cache cache.Store) Repository {
	return &postgresRepo{store: store, cache: cache}
}

func (r *postgresRepo) GetAll(ctx context.Context) ([]Parameter, error) {
	cacheKey := "parameter:all"
	var cached []Parameter
	found, err := r.cache.GetJSON(ctx, cacheKey, &cached)
	if err == nil && found {
		return cached, nil
	}

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

	_ = r.cache.SetJSON(ctx, cacheKey, result, 0) // Unlimited TTL
	return result, nil
}

func (r *postgresRepo) GetByID(ctx context.Context, id int) (Parameter, error) {
	cacheKey := fmt.Sprintf("parameter:%d", id)
	var cached Parameter
	found, err := r.cache.GetJSON(ctx, cacheKey, &cached)
	if err == nil && found {
		return cached, nil
	}

	sql := `SELECT id, name, unit, parameter_group, description, created_at, updated_at FROM parameter WHERE id = $1`
	rows, err := r.store.ExecQuery(ctx, sql, id)
	if err != nil || len(rows) == 0 {
		return Parameter{}, err
	}

	row := rows[0]
	p := Parameter{
		ID:             int(row["id"].(int32)),
		Name:           row["name"].(string),
		Unit:           row["unit"].(string),
		ParameterGroup: row["parameter_group"].(string),
		Description:    row["description"].(string),
		CreatedAt:      row["created_at"].(time.Time),
		UpdatedAt:      row["updated_at"].(time.Time),
	}

	_ = r.cache.SetJSON(ctx, cacheKey, p, 0) // Unlimited TTL
	return p, nil
}

func (r *postgresRepo) Create(ctx context.Context, p Parameter) error {
	sql := `INSERT INTO parameter(name, unit, parameter_group, description, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`
	err := r.store.Exec(ctx, sql, p.Name, p.Unit, p.ParameterGroup, p.Description, p.CreatedAt, p.UpdatedAt)
	if err == nil {
		_ = r.cache.Delete(ctx, "parameter:all") // clear list cache
	}
	return err
}

func (r *postgresRepo) Update(ctx context.Context, p Parameter) error {
	sql := `UPDATE parameter 
			SET name = $1, unit = $2, parameter_group = $3, description = $4, updated_at = $5 
			WHERE id = $6`
	err := r.store.Exec(ctx, sql, p.Name, p.Unit, p.ParameterGroup, p.Description, p.UpdatedAt, p.ID)
	if err == nil {
		_ = r.cache.Delete(ctx, fmt.Sprintf("parameter:%d", p.ID))
		_ = r.cache.Delete(ctx, "parameter:all")
	}
	return err
}

func (r *postgresRepo) Delete(ctx context.Context, id int) error {
	sql := `DELETE FROM parameter WHERE id = $1`
	err := r.store.Exec(ctx, sql, id)
	if err == nil {
		_ = r.cache.Delete(ctx, fmt.Sprintf("parameter:%d", id))
		_ = r.cache.Delete(ctx, "parameter:all")
	}
	return err
}
