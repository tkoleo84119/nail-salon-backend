package utils

import (
	"strings"
)

// HandleSortByMap handles the sort parameter by map => for join table sort
func HandleSortByMap(allowedFields map[string]string, defaultSortArr []string, sort *[]string) string {
	var sortParts []string

	if sort != nil && len(*sort) > 0 {
		for _, s := range *sort {
			// if start with -, order by desc
			rowField := strings.TrimPrefix(s, "-")
			column, ok := allowedFields[rowField]
			if !ok {
				continue
			}

			if strings.HasPrefix(s, "-") {
				sortParts = append(sortParts, column+" DESC")
			} else {
				sortParts = append(sortParts, column+" ASC")
			}
		}
	}

	if len(sortParts) == 0 {
		sortParts = append(sortParts, defaultSortArr...)
	}

	return strings.Join(sortParts, ", ")
}
