package utils

import "fmt"

func GetFiltersFromQueryStrings(query map[string][]string, filterKeys []string) []string {
	filters := make([]string, 0)
	for _, value := range filterKeys {
		if query[value] != nil {
			var filter string
			switch value {
			case "id":
				filter = fmt.Sprintf("id = %s", query[value][0])
			case "name":
				filter = fmt.Sprintf("name LIKE '%%%s%%'", query[value][0])
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
