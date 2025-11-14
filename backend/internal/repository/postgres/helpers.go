package postgres

import "strings"

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "unique") ||
		strings.Contains(err.Error(), "duplicate")
}

func containsField(err error, field string) bool {
	return strings.Contains(strings.ToLower(err.Error()), field)
}
