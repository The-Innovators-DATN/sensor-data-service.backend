package repository

import (
	"context"
	"fmt"
	"log"

	"sensor-data-service.backend/internal/common/castutil"
	"sensor-data-service.backend/internal/domain/model"
	"sensor-data-service.backend/internal/infrastructure/cache"
	"sensor-data-service.backend/internal/infrastructure/db"
)

type DashboardDataRepository struct {
	store db.Store
	cache cache.Store
}

func NewDashboardDataRepository(store db.Store, cache cache.Store) *DashboardDataRepository {
	return &DashboardDataRepository{store: store, cache: cache}
}

func (r *DashboardDataRepository) FindByID(ctx context.Context, uid string, userID int32) (*model.Dashboard, error) {
	log.Printf("[repo] FindByID: uid=%s user_id=%d", uid, userID)

	// query := `
	// 	SELECT uid, name, description, layout_configuration, created_by, created_at, updated_at, version, status
	// 	FROM dashboard
	// 	WHERE uid = $1 AND created_by = $2`

	query := `
		SELECT uid, name, description, layout_configuration, created_by, created_at, updated_at, version, status
		FROM dashboard
		WHERE uid = $1`
	rows, err := r.store.ExecQuery(ctx, query, uid)
	if err != nil {
		log.Printf("[repo][error] FindByID query failed: %v", err)
		return nil, err
	}
	if len(rows) == 0 {
		log.Printf("[repo][warn] No dashboard found for uid=%s user_id=%d", uid, userID)
		return nil, fmt.Errorf("dashboard not found")
	}

	d := mapRowToDashboard(rows[0])
	return d, nil
}

func (r *DashboardDataRepository) FindByIDAndUser(ctx context.Context, uid string, userID int32) (*model.Dashboard, error) {
	log.Printf("[repo] FindByIDAndUser: uid=%s user=%d", uid, userID)
	// query := `SELECT uid, name, description, layout_configuration, created_by, created_at, updated_at, version, status FROM dashboard WHERE uid = $1 AND created_by = $2`
	query := `SELECT uid, name, description, layout_configuration, created_by, created_at, updated_at, version, status FROM dashboard WHERE uid = $1`
	rows, err := r.store.ExecQuery(ctx, query, uid)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	return mapRowToDashboard(rows[0]), nil
}

func (r *DashboardDataRepository) FindByNameAndUser(ctx context.Context, name string, userID int32) (*model.Dashboard, error) {
	log.Printf("[repo] FindByNameAndUser: name=%s user=%d", name, userID)
	query := `SELECT uid, name, description, layout_configuration, created_by, created_at, updated_at, version, status FROM dashboard WHERE name = $1 AND created_by = $2`
	rows, err := r.store.ExecQuery(ctx, query, name, userID)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	return mapRowToDashboard(rows[0]), nil
}

func (r *DashboardDataRepository) ListByUser(ctx context.Context, userID int32, offset, limit int32) ([]*model.Dashboard, error) {
	log.Printf("[repo] ListByUser: user=%d", userID)
	query := `
		SELECT uid, name, description, layout_configuration, created_by,
		       created_at, updated_at, version, status
		FROM dashboard
		WHERE created_by = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.store.ExecQuery(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query dashboards paged: %w", err)
	}
	if len(rows) == 0 {
		log.Printf("[repo][warn] No dashboards found for user_id=%d", userID)
		return nil, fmt.Errorf("no dashboards found")
	}
	return mapRowsToDashboards(rows), nil
}
func (r *DashboardDataRepository) CountDashboardsByUser(ctx context.Context, createdBy int32) (int32, error) {
	query := `SELECT COUNT(*) AS total FROM dashboard WHERE created_by = $1`

	rows, err := r.store.ExecQuery(ctx, query, createdBy)
	if err != nil {
		return 0, fmt.Errorf("query dashboard count: %w", err)
	}
	if len(rows) == 0 {
		return 0, fmt.Errorf("no result returned for count")
	}

	totalVal := rows[0]["total"]
	total, ok := totalVal.(int64) // pgx returns int64 for COUNT(*)
	if !ok {
		return 0, fmt.Errorf("invalid type for count: %T", totalVal)
	}
	return int32(total), nil
}

func (r *DashboardDataRepository) ListAll(ctx context.Context) ([]*model.Dashboard, error) {
	log.Printf("[repo] ListAll")
	query := `SELECT uid, name, description, layout_configuration, created_by, created_at, updated_at, version, status FROM dashboard ORDER BY uid DESC`
	rows, err := r.store.ExecQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	return mapRowsToDashboards(rows), nil
}

func (r *DashboardDataRepository) Create(ctx context.Context, d *model.Dashboard) (string, error) {
	log.Printf("[repo] Create: uid=%s user_id=%d", d.UID, d.CreatedBy)
	query := `
		INSERT INTO dashboard (uid, name, description, layout_configuration, created_by, created_at, updated_at, version, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING uid
	`
	// var uid string
	rows, err := r.store.ExecQuery(ctx, query,
		d.UID,
		d.Name,
		d.Description,
		d.LayoutConfiguration,
		d.CreatedBy,
		d.CreatedAt,
		d.UpdatedAt,
		d.Version,
		d.Status,
	)
	if err != nil {
		log.Printf("[repo][error] Create query failed: %v", err)
		return "", err
	}
	if len(rows) == 0 {
		log.Printf("[repo][warn] No uid returned for dashboard creation")
		return "", fmt.Errorf("no uid returned for dashboard creation")
	}
	uidVal := castutil.ToUUID(rows[0]["uid"])
	// uid, ok := uidVal.(string)
	// if !ok {
	// 	log.Printf("[repo][error] Invalid type for uid: %T", uidVal)
	// 	return "", fmt.Errorf("invalid type for uid: %T", uidVal)
	// }
	return uidVal.String(), nil
}

func (r *DashboardDataRepository) Update(ctx context.Context, d *model.Dashboard) error {
	log.Printf("Raw dashboard: %v", d)
	// log.Printf(d.LayoutConfiguration)
	query := `
		UPDATE dashboard
		SET name=$1, description=$2, layout_configuration=$3, status=$4
		WHERE uid = $5`
	// print type of d.LayoutConfiguration
	err := r.store.Exec(ctx, query, d.Name, d.Description, d.LayoutConfiguration, d.Status, d.UID)
	if err != nil {
		log.Printf("[repo][error] Update query failed: %v", err)
	}

	return err

}

func (r *DashboardDataRepository) Patch(ctx context.Context, d *model.Dashboard) error {
	log.Printf("[repo] Patch: uid=%s user_id=%d", d.UID, d.CreatedBy)
	query := `
		UPDATE dashboard
		SET layout_configuration=$1, updated_at=NOW(), version=$2
		WHERE uid = $3`
	return r.store.Exec(ctx, query, d.LayoutConfiguration, d.Version, d.UID)
}

func (r *DashboardDataRepository) Delete(ctx context.Context, uid string) error {
	log.Printf("[repo] Delete: uid=%s", uid)
	query := `DELETE FROM dashboard WHERE uid = $1`
	return r.store.Exec(ctx, query, uid)
}

func mapRowToDashboard(r map[string]interface{}) *model.Dashboard {

	return &model.Dashboard{
		UID:                 castutil.ToUUID(r["uid"]),
		Name:                castutil.ToString(r["name"]),
		Description:         castutil.ToString(r["description"]),
		LayoutConfiguration: castutil.ToString(r["layout_configuration"]),
		CreatedBy:           int32(castutil.ToInt(r["created_by"])),
		CreatedAt:           castutil.ToTime(r["created_at"]),
		UpdatedAt:           castutil.ToTime(r["updated_at"]),
		Version:             int32(castutil.ToInt(r["version"])),
		Status:              castutil.ToString(r["status"]),
	}
}

func mapRowsToDashboards(rows []map[string]interface{}) []*model.Dashboard {
	var list []*model.Dashboard
	for _, r := range rows {
		list = append(list, mapRowToDashboard(r))
	}
	return list
}
