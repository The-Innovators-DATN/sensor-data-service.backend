package parameter

import "time"

func toString(val interface{}) string {
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

func toInt(val interface{}) int {
	if i, ok := val.(int32); ok {
		return int(i)
	}
	if i, ok := val.(int64); ok {
		return int(i)
	}
	return 0
}

func toTime(val interface{}) time.Time {
	if t, ok := val.(time.Time); ok {
		return t
	}
	return time.Time{}
}
