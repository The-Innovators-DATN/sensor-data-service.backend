package castutil

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

func ToInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}

func ToFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

func ToBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case uint8, int, int64, float64:
		return ToInt(val) != 0
	case string:
		return val == "1" || val == "true"
	default:
		return false
	}
}

func ToString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func ToTime(v interface{}) time.Time {
	if t, ok := v.(time.Time); ok {
		return t
	}
	return time.Time{}
}

func MustToFloat(v interface{}) float64 {
	f, _ := ToFloat(v)
	return f
}
func TimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func Float32PointerOrZero(f *float64) float32 {
	if f == nil {
		return 0
	}
	return float32(*f)
}
func OptionalFloat(v interface{}) *float64 {
	if v == nil {
		return nil
	}
	switch f := v.(type) {
	case float64:
		return &f
	case float32:
		val := float64(f)
		return &val
	case int:
		val := float64(f)
		return &val
	case int64:
		val := float64(f)
		return &val
	default:
		return nil
	}
}

func OptionalTime(v interface{}) *time.Time {
	if v == nil {
		return nil
	}
	if t, ok := v.(time.Time); ok {
		return &t
	}
	return nil
}

func FormatOptionalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
func ToInt32Array(ids []int32) pgtype.Int4Array {
	arr := pgtype.Int4Array{}
	_ = arr.Set(ids) // convert slice to pg-compatible array
	return arr
}

func ParseUUID(val interface{}) (uuid.UUID, error) {
	switch v := val.(type) {
	case string:
		return uuid.Parse(v)
	case []byte:
		return uuid.FromBytes(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, fmt.Errorf("unexpected type for UUID: %T", val)
	}
}

// Không trả error, panic nếu có lỗi (cho nhanh gọn)
func MustParseUUID(val interface{}) uuid.UUID {
	u, err := ParseUUID(val)
	if err != nil {
		panic(err)
	}
	return u
}

func ToUUID(val interface{}) uuid.UUID {
	switch v := val.(type) {
	case uuid.UUID:
		return v
	case [16]byte:
		return uuid.UUID(v)
	case []byte:
		u, err := uuid.FromBytes(v)
		if err == nil {
			return u
		}
	case string:
		u, err := uuid.Parse(v)
		if err == nil {
			return u
		}
	}
	log.Printf("ToUUID: unsupported type or invalid UUID: %T → returning uuid.Nil", val)
	return uuid.Nil
}
