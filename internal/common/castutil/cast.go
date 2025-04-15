package castutil

import "time"

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
