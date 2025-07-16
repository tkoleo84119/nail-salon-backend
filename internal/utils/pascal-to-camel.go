package utils

import (
	"strings"
	"unicode"
)

func PascalToCamel(s string) string {
	if s == "" {
		return s
	}

	if strings.ToUpper(s) == s {
		return strings.ToLower(s)
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	for i := 1; i < len(r); {
		if unicode.IsUpper(r[i]) {
			start := i
			for i < len(r) && unicode.IsUpper(r[i]) {
				i++
			}
			if i-start > 1 {
				for j := start + 1; j < i; j++ {
					r[j] = unicode.ToLower(r[j])
				}
			}
		} else {
			i++
		}
	}
	return string(r)
}
