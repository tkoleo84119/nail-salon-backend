package utils

import (
	"slices"
	"strings"
)

// HandleSort handles the sort parameter
func HandleSort(allowedFields []string, defaultSort string, defaultSortDirection string, sort *[]string) string {
	sortStr := ""
	if sort != nil && len(*sort) > 0 {
		for _, s := range *sort {
			// if start with -, order by desc
			field := CamelToSnake(strings.TrimPrefix(s, "-"))
			if !slices.Contains(allowedFields, field) {
				continue
			}

			if strings.HasPrefix(s, "-") {
				sortStr = field + " DESC"
			} else {
				sortStr = field + " ASC"
			}
		}
	}
	if sortStr == "" {
		sortStr = defaultSort + " " + defaultSortDirection
	}

	return sortStr
}
