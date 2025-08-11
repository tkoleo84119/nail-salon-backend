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

// HandleSortByMap handles the sort parameter by map => for join table sort
func HandleSortByMap(allowedFields map[string]string, defaultSort string, defaultSortDirection string, sort *[]string) string {
	sortStr := ""
	if sort != nil && len(*sort) > 0 {
		for _, s := range *sort {
			// if start with -, order by desc
			field := CamelToSnake(strings.TrimPrefix(s, "-"))
			column, ok := allowedFields[field]
			if !ok {
				continue
			}

			if strings.HasPrefix(s, "-") {
				sortStr = column + " DESC"
			} else {
				sortStr = column + " ASC"
			}
		}
	}
	if sortStr == "" {
		sortStr = defaultSort + " " + defaultSortDirection
	}
	return sortStr
}
