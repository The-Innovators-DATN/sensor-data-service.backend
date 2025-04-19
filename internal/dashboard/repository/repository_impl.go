package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"sensor-data-service.backend/infrastructure/cache"
	"sensor-data-service.backend/infrastructure/db"
	"sensor-data-service.backend/internal/common/castutil"
	"sensor-data-service.backend/internal/dashboard/domain"
)

type DashboardDataRepository struct {
	store db.Store
	cache cache.Store
}

func NewDashboardDataRepository(store db.Store, cache cache.Store) *DashboardDataRepository {
	return &DashboardDataRepository{store: store, cache: cache}
}

func (r *DashboardDataRepository) FindByID(ctx context.Context, id int32) (*domain.Dashboard, error) {
	cacheKey := fmt.Sprintf("dashboard:%d", id)
	var cached domain.Dashboard
	found, err := r.cache.GetJSON(ctx, cacheKey, &cached)
	if err != nil {
		log.Printf("[warn] FindByID cache get error: %v", err)
	}
	if found {
		return &cached, nil
	}

	query := `SELECT id, name, description, layout_json, created_by, created_at, updated_at, version, status FROM dashboard WHERE id = $1`
	rows, err := r.store.ExecQuery(ctx, query, id)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	rw := rows[0]
	d := &domain.Dashboard{
		ID:          int32(castutil.ToInt(rw["id"])),
		Name:        castutil.ToString(rw["name"]),
		Description: castutil.ToString(rw["description"]),
		LayoutJSON:  castutil.ToString(rw["layout_json"]),
		CreatedBy:   int32(castutil.ToInt(rw["created_by"])),
		CreatedAt:   castutil.ToString(rw["created_at"]),
		UpdatedAt:   castutil.ToString(rw["updated_at"]),
		Version:     int32(castutil.ToInt(rw["version"])),
		Status:      castutil.ToString(rw["status"]),
	}

	_ = r.cache.SetJSON(ctx, cacheKey, d, int64(time.Hour.Seconds()))
	return d, nil
}

func (r *DashboardDataRepository) List(ctx context.Context) ([]*domain.Dashboard, error) {
	query := `SELECT id, name, description, layout_json, created_by, created_at, updated_at, version, status FROM dashboard ORDER BY id DESC`
	rows, err := r.store.ExecQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	var dashboards []*domain.Dashboard
	for _, rw := range rows {
		dashboards = append(dashboards, &domain.Dashboard{
			ID:          int32(castutil.ToInt(rw["id"])),
			Name:        castutil.ToString(rw["name"]),
			Description: castutil.ToString(rw["description"]),
			LayoutJSON:  castutil.ToString(rw["layout_json"]),
			CreatedBy:   int32(castutil.ToInt(rw["created_by"])),
			CreatedAt:   castutil.ToString(rw["created_at"]),
			UpdatedAt:   castutil.ToString(rw["updated_at"]),
			Version:     int32(castutil.ToInt(rw["version"])),
			Status:      castutil.ToString(rw["status"]),
		})
	}
	return dashboards, nil
}

func (r *DashboardDataRepository) Save(ctx context.Context, d *domain.Dashboard) error {
	query := `INSERT INTO dashboard (name, description, layout_json, created_by, created_at, updated_at, version, status)
	VALUES ($1, $2, $3, $4, NOW(), NOW(), $5, $6)
	ON CONFLICT (id) DO UPDATE SET name=$1, description=$2, layout_json=$3, updated_at=NOW(), version=$5, status=$6`
	return r.store.Exec(ctx, query, d.Name, d.Description, d.LayoutJSON, d.CreatedBy, d.Version, d.Status)
}

func (r *DashboardDataRepository) Delete(ctx context.Context, id int32) error {
	cacheKey := fmt.Sprintf("dashboard:%d", id)
	_ = r.cache.Delete(ctx, cacheKey)
	query := `DELETE FROM dashboard WHERE id = $1`
	return r.store.Exec(ctx, query, id)
}

func (r *DashboardDataRepository) Update(ctx context.Context, d *domain.Dashboard) error {
	cacheKey := fmt.Sprintf("dashboard:%d", d.ID)
	_ = r.cache.Delete(ctx, cacheKey)
	query := `UPDATE dashboard SET name=$1, description=$2, layout_json=$3, updated_at=NOW(), version=$4, status=$5 WHERE id = $6`
	return r.store.Exec(ctx, query, d.Name, d.Description, d.LayoutJSON, d.Version, d.Status, d.ID)
}

func (r *DashboardDataRepository) FindByName(ctx context.Context, name string) (*domain.Dashboard, error) {
	cacheKey := fmt.Sprintf("dashboard:name:%s", name)
	var cached domain.Dashboard
	found, err := r.cache.GetJSON(ctx, cacheKey, &cached)
	if err != nil {
		log.Printf("[warn] FindByName cache get error: %v", err)
	}
	if found {
		return &cached, nil
	}

	query := `SELECT id, name, description, layout_json, created_by, created_at, updated_at, version, status FROM dashboard WHERE name = $1`
	rows, err := r.store.ExecQuery(ctx, query, name)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	rw := rows[0]
	d := &domain.Dashboard{
		ID:          int32(castutil.ToInt(rw["id"])),
		Name:        castutil.ToString(rw["name"]),
		Description: castutil.ToString(rw["description"]),
		LayoutJSON:  castutil.ToString(rw["layout_json"]),
		CreatedBy:   int32(castutil.ToInt(rw["created_by"])),
		CreatedAt:   castutil.ToString(rw["created_at"]),
		UpdatedAt:   castutil.ToString(rw["updated_at"]),
		Version:     int32(castutil.ToInt(rw["version"])),
		Status:      castutil.ToString(rw["status"]),
	}

	_ = r.cache.SetJSON(ctx, cacheKey, d, int64(time.Hour.Seconds()))
	return d, nil
}
