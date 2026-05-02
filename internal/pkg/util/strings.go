package util

import (
	"strings"
)

func StringPointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func NonEmptyStrings(values ...string) []string {
	nonEmpty := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			nonEmpty = append(nonEmpty, value)
		}
	}
	return nonEmpty
}
