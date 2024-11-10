package helpers

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

type Filter struct {
	Field     string
	Operation string
	Value     interface{}
}

func applyFilters(query squirrel.SelectBuilder, filters []Filter) (squirrel.SelectBuilder, error) {
	for _, filter := range filters {
		switch filter.Operation {
		case "=":
			query = query.Where(squirrel.Eq{filter.Field: filter.Value})
		case ">":
			query = query.Where(squirrel.Gt{filter.Field: filter.Value})
		case "<":
			query = query.Where(squirrel.Lt{filter.Field: filter.Value})
		case ">=":
			query = query.Where(squirrel.GtOrEq{filter.Field: filter.Value})
		case "<=":
			query = query.Where(squirrel.LtOrEq{filter.Field: filter.Value})
		case "LIKE":
			query = query.Where(squirrel.Like{filter.Field: filter.Value})
		case "!=":
			query = query.Where(squirrel.NotEq{filter.Field: filter.Value})
		case "IN":
			if vals, ok := filter.Value.([]interface{}); ok {
				query = query.Where(squirrel.Eq{filter.Field: vals})
			} else {
				return query, fmt.Errorf("invalid value type for IN filter on field %s", filter.Field)
			}
		case "NOT IN":
			if vals, ok := filter.Value.([]interface{}); ok {
				query = query.Where(squirrel.NotEq{filter.Field: vals})
			} else {
				return query, fmt.Errorf("invalid value type for NOT IN filter on field %s", filter.Field)
			}
		// Add other cases as needed
		default:
			return query, fmt.Errorf("unsupported operation: %s", filter.Operation)
		}
	}
	return query, nil
}
