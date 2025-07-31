package utils

import (
	"strings"
	"unicode"
)

// CamelToSnake converts a camel case string to a snake case string
func CamelToSnake(s string) string {
	if s == "" {
		return s
	}

	var builder strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(s[i-1])
				if unicode.IsLower(prev) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) {
					builder.WriteByte('_')
				}
			}
			builder.WriteRune(unicode.ToLower(r))
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
