package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/infrastructure/cache"
	"sensor-data-service.backend/infrastructure/db"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/common/castutil"
	"sensor-data-service.backend/internal/station/domain"
)

type StationDataRepository struct {
	store db.Store
	cache cache.Store
}

func NewStationDataRepository(store db.Store, cache cache.Store) *StationDataRepository {
	return &StationDataRepository{
		store: store,
		cache: cache,
	}
}

// ========== BASIC CRUD ==========

func (r *StationDataRepository) FindStationByID(ctx context.Context, id int32) (*domain.StationWithLocation, error) {
	log.Printf("[debug] FindStationByID called with id: %d", id)

	query := `
	SELECT 
		s.id, s.name, s.description, s.lat, s.long, s.status, s.station_type, s.country, 
		s.water_body_id, s.station_manager, s.created_at, s.updated_at,
		w.name AS water_body_name, w.type AS water_body_type,
		c.id AS catchment_id, c.name AS catchment_name, c.description AS catchment_description,
		rb.id AS river_basin_id, rb.name AS river_basin_name
	FROM station s
	JOIN water_body w ON s.water_body_id = w.id
	JOIN catchment c ON w.catchment_id = c.id
	JOIN river_basin rb ON c.river_basin_id = rb.id
	WHERE s.id = $1 AND s.status != 'deleted'
	`

	rows, err := r.store.ExecQuery(ctx, query, id)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	row := rows[0]

	station := domain.Station{
		ID:             int32(castutil.ToInt(row["id"])),
		Name:           castutil.ToString(row["name"]),
		Description:    castutil.ToString(row["description"]),
		Lat:            float32(castutil.MustToFloat(row["lat"])),
		Long:           float32(castutil.MustToFloat(row["long"])),
		Status:         castutil.ToString(row["status"]),
		StationType:    castutil.ToString(row["station_type"]),
		Country:        castutil.ToString(row["country"]),
		WaterBodyID:    int32(castutil.ToInt(row["water_body_id"])),
		StationManager: int32(castutil.ToInt(row["station_manager"])),
		CreatedAt:      castutil.ToTime(row["created_at"]),
		UpdatedAt:      castutil.ToTime(row["updated_at"]),
	}

	location := domain.StationLocation{
		WaterBodyName:  castutil.ToString(row["water_body_name"]),
		WaterBodyType:  castutil.ToString(row["water_body_type"]),
		CatchmentID:    int32(castutil.ToInt(row["catchment_id"])),
		CatchmentName:  castutil.ToString(row["catchment_name"]),
		CatchmentDesc:  castutil.ToString(row["catchment_description"]),
		RiverBasinID:   int32(castutil.ToInt(row["river_basin_id"])),
		RiverBasinName: castutil.ToString(row["river_basin_name"]),
	}

	return &domain.StationWithLocation{
		Station:  station,
		Location: location,
	}, nil
}

func (r *StationDataRepository) FindStationIDsByRiverBasin(ctx context.Context, riverBasinID int32) ([]int32, error) {
	query := `
		SELECT s.id
		FROM station s
		JOIN water_body w ON s.water_body_id = w.id
		JOIN catchment c ON w.catchment_id = c.id
		WHERE c.river_basin_id = $1 AND s.status != 'deleted'
	`
	rows, err := r.store.ExecQuery(ctx, query, riverBasinID)
	if err != nil {
		return nil, err
	}

	var ids []int32
	for _, row := range rows {
		ids = append(ids, int32(castutil.ToInt(row["id"])))
	}
	return ids, nil
}

func (r *StationDataRepository) FindStationIDsByCatchment(ctx context.Context, catchmentID int32) ([]int32, error) {
	log.Printf("[debug] FindStationIDsByCatchment called with catchmentID: %d", catchmentID)
	query := `
		SELECT s.id
		FROM station s
		JOIN water_body w ON s.water_body_id = w.id
		WHERE w.catchment_id = $1 AND s.status != 'deleted'
	`
	rows, err := r.store.ExecQuery(ctx, query, catchmentID)
	if err != nil {
		return nil, err
	}

	var ids []int32
	for _, row := range rows {
		ids = append(ids, int32(castutil.ToInt(row["id"])))
	}
	return ids, nil
}
func (r *StationDataRepository) FindStationIDsByWaterBody(ctx context.Context, waterBodyID int32) ([]int32, error) {
	log.Printf("[debug] FindStationIDsByWaterBody called with waterBodyID: %d", waterBodyID)
	query := `SELECT id FROM station WHERE water_body_id = $1 AND status != 'deleted'`
	rows, err := r.store.ExecQuery(ctx, query, waterBodyID)
	if err != nil {
		return nil, err
	}

	var ids []int32
	for _, row := range rows {
		ids = append(ids, int32(castutil.ToInt(row["id"])))
	}
	return ids, nil
}

func (r *StationDataRepository) FilterStations(ctx context.Context, keyword string, waterBodyID, catchmentID, riverBasinID int32) ([]*domain.Station, error) {
	log.Printf("[debug] FilterStations called with keyword: %s, waterBodyID: %d, catchmentID: %d, riverBasinID: %d", keyword, waterBodyID, catchmentID, riverBasinID)
	query := `
		SELECT s.id, s.name, s.description, s.lat, s.long, s.status, s.station_type, s.country, s.water_body_id, s.station_manager, s.created_at, s.updated_at
		FROM station s
		JOIN water_body w ON s.water_body_id = w.id
		JOIN catchment c ON w.catchment_id = c.id
		WHERE s.status != 'deleted'
	`
	args := []interface{}{}
	if keyword != "" {
		query += " AND s.name ILIKE '%' || $1 || '%'"
		args = append(args, keyword)
	}
	if waterBodyID > 0 {
		query += fmt.Sprintf(" AND s.water_body_id = $%d", len(args)+1)
		args = append(args, waterBodyID)
	}
	if catchmentID > 0 {
		query += fmt.Sprintf(" AND c.id = $%d", len(args)+1)
		args = append(args, catchmentID)
	}
	if riverBasinID > 0 {
		query += fmt.Sprintf(" AND c.river_basin_id = $%d", len(args)+1)
		args = append(args, fmt.Sprint(riverBasinID))
	}

	rows, err := r.store.ExecQuery(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var stations []*domain.Station
	for _, row := range rows {
		stations = append(stations, &domain.Station{
			ID:             int32(castutil.ToInt(row["id"])),
			Name:           castutil.ToString(row["name"]),
			Description:    castutil.ToString(row["description"]),
			Lat:            float32(castutil.MustToFloat(row["lat"])),
			Long:           float32(castutil.MustToFloat(row["long"])),
			Status:         castutil.ToString(row["status"]),
			StationType:    castutil.ToString(row["station_type"]),
			Country:        castutil.ToString(row["country"]),
			WaterBodyID:    int32(castutil.ToInt(row["water_body_id"])),
			StationManager: int32(castutil.ToInt(row["station_manager"])),
			CreatedAt:      castutil.ToTime(row["created_at"]),
			UpdatedAt:      castutil.ToTime(row["updated_at"]),
		})
	}
	return stations, nil
}

func (r *StationDataRepository) InsertStation(ctx context.Context, st domain.Station) error {
	log.Printf("[debug] InsertStation called with station: %+v", st)
	query := `
		INSERT INTO station 
		(name, description, lat, long, status, station_type, country, water_body_id, station_manager, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`
	err := r.store.Exec(ctx, query,
		st.Name, st.Description, st.Lat, st.Long, st.Status,
		st.StationType, st.Country, st.WaterBodyID, st.StationManager)
	return err
}

func (r *StationDataRepository) UpdateStationStatus(ctx context.Context, id int32, status string) error {
	log.Printf("[debug] UpdateStationStatus called with id: %d, status: %s", id, status)
	query := `UPDATE station SET status = $1, updated_at = NOW() WHERE id = $2`
	err := r.store.Exec(ctx, query, status, id)
	return err
}

// ========== CACHED TARGET RESOLUTION ==========

func (r *StationDataRepository) GetStationsByTarget(ctx context.Context, targetType stationpb.TargetType, targetId int32) ([]int32, error) {

	log.Print("[debug] GetStationsByTarget called")
	cacheKey := fmt.Sprintf("station_targets:%d:%d", targetType, targetId)

	var ids []int32
	found, err := r.cache.GetJSON(ctx, cacheKey, &ids)
	if err != nil {
		log.Printf("[warn] Redis error: %v", err)
	}
	if found {
		return ids, nil
	}

	var sql string
	switch targetType {
	case stationpb.TargetType_STATION:
		sql = "SELECT id FROM station WHERE id = $1"
	case stationpb.TargetType_WATER_BODY:
		sql = "SELECT id FROM station WHERE water_body_id = $1"
	case stationpb.TargetType_CATCHMENT:
		sql = `
			SELECT s.id FROM station s
			JOIN water_body w ON s.water_body_id = w.id
			WHERE w.catchment_id = $1
		`
	case stationpb.TargetType_RIVER_BASIN:
		sql = `
			SELECT s.id FROM station s
			JOIN water_body w ON s.water_body_id = w.id
			JOIN catchment c ON w.catchment_id = c.id
			WHERE c.river_basin_id = $1
		`
	default:
		return nil, fmt.Errorf("unsupported targetType: %v", targetType)
	}

	rows, err := r.store.ExecQuery(ctx, sql, targetId)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		ids = append(ids, int32(castutil.ToInt(row["id"])))
	}

	_ = r.cache.SetJSON(ctx, cacheKey, ids, int64(time.Hour.Seconds()))
	return ids, nil
}

// ========== ENUM VALUE LOADER ==========

func (r *StationDataRepository) ListEnumValues(table string) ([]*common.EnumValue, error) {
	query := fmt.Sprintf("SELECT name FROM %s", table)
	rows, err := r.store.ExecQuery(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var result []*common.EnumValue
	for _, row := range rows {
		result = append(result, &common.EnumValue{Name: castutil.ToString(row["name"])})
	}
	return result, nil
}
