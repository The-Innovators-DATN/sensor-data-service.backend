package domain

import (
	"time"
)

// Station tương ứng với bảng "station"
type Station struct {
	ID             int       `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	Status         string    `db:"status"`
	Long           float32   `db:"long"`
	Lat            float32   `db:"lat"`
	Country        string    `db:"country"`
	StationType    string    `db:"station_type"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	StationManager int       `db:"station_manager"`
	WaterBodyID    int       `db:"water_body_id"`
}

// WaterBody tương ứng với bảng "water_body"
type WaterBody struct {
	ID          int       `db:"id"`
	Type        string    `db:"type"`
	Name        string    `db:"name"`
	CatchmentID int       `db:"catchment_id"`
	UpdatedAt   time.Time `db:"updated_at"`
	Description string    `db:"description"`
}

// Parameter tương ứng với bảng "parameter"
type Parameter struct {
	ID             int       `db:"id"`
	Name           string    `db:"name"`
	Unit           string    `db:"unit"`
	ParameterGroup *string   `db:"parameter_group"` // Có thể null
	Description    *string   `db:"description"`     // Có thể null
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// StationParameter tương ứng với bảng "station_parameter"
// Primary key là (parameter_id, station_id)
type StationParameter struct {
	ParameterID  int        `db:"parameter_id"`
	StationID    int        `db:"station_id"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	LastReceivAt *time.Time `db:"last_receiv_at"` // Có thể null nếu chưa nhận giá trị
	LastValue    *float32   `db:"last_value"`     // Có thể null nếu chưa có giá trị
}

// Catchment tương ứng với bảng "catchment"
// Lưu ý: "river_basin_id" của bảng này định nghĩa kiểu VARCHAR(255) nên sử dụng string.
type Catchment struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	RiverBasinID string    `db:"river_basin_id"`
	Country      string    `db:"country"`
	Description  *string   `db:"description"` // Có thể null
	UpdatedAt    time.Time `db:"updated_at"`
}

// RiverBasin tương ứng với bảng "river_basin"
type RiverBasin struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Country tương ứng với bảng "country"
type Country struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	UpdatedAt time.Time `db:"updated_at"`
}

// StarDashboard tương ứng với bảng "star_dashboard"
// Với cột "id" kiểu UUID, ta sử dụng kiểu string để lưu trữ dưới dạng chuỗi.
type StarDashboard struct {
	ID                  string    `db:"id"`
	UserID              int       `db:"user_id"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"update_at"` // Chú ý: cột tên "update_at" theo định nghĩa DB
	Version             int       `db:"version"`
	LayoutConfiguration string    `db:"layout_configuration"`
}

// StationKey tương ứng với bảng "station_key"
type StationKey struct {
	ID        int `db:"id"`
	StationID int `db:"station_id"`
	OrgID     int `db:"org_id"`
	// is_revoked kiểu SMALLINT, có thể dùng int hoặc bool (nếu xử lý chuyển đổi)
	IsRevoked int       `db:"is_revoked"`
	Name      string    `db:"name"`
	Key       string    `db:"key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"update_at"`
}

// StationUploadStatus tương ứng với bảng "station_upload_status"
type StationUploadStatus struct {
	UID          int       `db:"uid"`
	AttachmentID int       `db:"attachment_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Status       string    `db:"status"`
	// user_id có thể là null, do đó sử dụng con trỏ hoặc sql.NullInt64
	UserID *int `db:"user_id"`
}

// StationAttachments tương ứng với bảng "station_attachments"
// Với uid kiểu UUID, sử dụng string để lưu trữ.
type StationAttachments struct {
	UID           string    `db:"uid"`
	Size          int       `db:"size"`
	Filename      string    `db:"filename"`
	ContentType   string    `db:"content_type"`
	DisplayName   string    `db:"display_name"`
	WorkflowState string    `db:"workflow_state"`
	UserID        int       `db:"user_id"`
	FileState     string    `db:"file_state"`
	Namespace     string    `db:"namespace"`
	StationID     int       `db:"station_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
