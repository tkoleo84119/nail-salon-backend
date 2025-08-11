package utils

import "strings"

func SetDefaultValuesOfPagination(limit *int, offset *int, limitDefault int, offsetDefault int) (int, int) {
	limitValue := limitDefault
	offsetValue := offsetDefault

	if limit != nil && *limit > 0 {
		limitValue = *limit
	}
	if offset != nil && *offset >= 0 {
		offsetValue = *offset
	}

	return limitValue, offsetValue
}

func TransformSort(sort *string) []string {
	if sort == nil || *sort == "" {
		return []string{}
	}

	return strings.Split(*sort, ",")
}
