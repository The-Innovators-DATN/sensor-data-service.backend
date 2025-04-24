package utils

// SanitizePagination ensures default values for pagination and clamps invalid inputs
func SanitizePagination(page, limit int32) (int32, int32) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return page, limit
}
