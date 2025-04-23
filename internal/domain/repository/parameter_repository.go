package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"sensor-data-service.backend/internal/domain/model"
	"sensor-data-service.backend/internal/infrastructure/cache"
	"sensor-data-service.backend/internal/infrastructure/db"
)

type ParameterRepository interface {
	GetAll(ctx context.Context) ([]model.Parameter, error)
	GetByID(ctx context.Context, id int) (model.Parameter, error)
	Create(ctx context.Context, p model.Parameter) error
	Update(ctx context.Context, p model.Parameter) error
	Delete(ctx context.Context, id int) error
}

type postgresRepo struct {
	store db.Store
	cache cache.Store
}

func NewParameterRepository(store db.Store, cache cache.Store) ParameterRepository {
	return &postgresRepo{store: store, cache: cache}
}

func (r *postgresRepo) GetAll(ctx context.Context) ([]model.Parameter, error) {
	log.Printf("[debug] GetAll called")
	cacheKey := "parameter:all"
	var cached []model.Parameter
	found, err := r.cache.GetJSON(ctx, cacheKey, &cached)
	if err == nil && found {
		return cached, nil
	}

	sql := `SELECT id, name, unit, parameter_group, description, created_at, updated_at FROM parameter`
	rows, err := r.store.ExecQuery(ctx, sql)
	if err != nil {
		return nil, err
	}

	var result []model.Parameter
	for _, row := range rows {
		p := model.Parameter{
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

func (r *postgresRepo) GetByID(ctx context.Context, id int) (model.Parameter, error) {
	log.Printf("[debug] GetByID called with id: %d", id)
	cacheKey := fmt.Sprintf("parameter:%d", id)
	var cached model.Parameter
	found, err := r.cache.GetJSON(ctx, cacheKey, &cached)
	if err == nil && found {
		return cached, nil
	}

	sql := `SELECT id, name, unit, parameter_group, description, created_at, updated_at FROM parameter WHERE id = $1`
	rows, err := r.store.ExecQuery(ctx, sql, id)
	if err != nil || len(rows) == 0 {
		return model.Parameter{}, err
	}

	row := rows[0]
	p := model.Parameter{
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

func (r *postgresRepo) Create(ctx context.Context, p model.Parameter) error {
	log.Printf("[debug] Create called with parameter: %+v", p)
	sql := `INSERT INTO parameter(name, unit, parameter_group, description, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`
	err := r.store.Exec(ctx, sql, p.Name, p.Unit, p.ParameterGroup, p.Description, p.CreatedAt, p.UpdatedAt)
	if err == nil {
		_ = r.cache.Delete(ctx, "parameter:all") // clear list cache
	}
	return err
}

func (r *postgresRepo) Update(ctx context.Context, p model.Parameter) error {
	log.Printf("[debug] Update called with parameter: %+v", p)
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
	log.Printf("[debug] Delete called with id: %d", id)
	sql := `DELETE FROM parameter WHERE id = $1`
	err := r.store.Exec(ctx, sql, id)
	if err == nil {
		_ = r.cache.Delete(ctx, fmt.Sprintf("parameter:%d", id))
		_ = r.cache.Delete(ctx, "parameter:all")
	}
	return err
}
