package helpers

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

type Filter struct {
	Field     FilterKey
	Operation string
	Value     interface{}
}

type FilterKey string

const (
	FilterKeyTag       FilterKey = "tag"
	FilterKeyCategogy  FilterKey = "category"
	FilterKeyCreatedAt FilterKey = "createdAt"
	FilterKeyCopiedAt  FilterKey = "copiedAt"
)

func applyFilters(query squirrel.SelectBuilder, filters []Filter) (squirrel.SelectBuilder, error) {
	for _, filter := range filters {
		switch filter.Operation {
		case "=":
			query = query.Where(squirrel.Eq{string(filter.Field): filter.Value})
		case ">":
			query = query.Where(squirrel.Gt{string(filter.Field): filter.Value})
		case "<":
			query = query.Where(squirrel.Lt{string(filter.Field): filter.Value})
		case ">=":
			query = query.Where(squirrel.GtOrEq{string(filter.Field): filter.Value})
		case "<=":
			query = query.Where(squirrel.LtOrEq{string(filter.Field): filter.Value})
		case "LIKE":
			query = query.Where(squirrel.Like{string(filter.Field): filter.Value})
		case "!=":
			query = query.Where(squirrel.NotEq{string(filter.Field): filter.Value})
		case "IN":
			if vals, ok := filter.Value.([]interface{}); ok {
				query = query.Where(squirrel.Eq{string(filter.Field): vals})
			} else {
				return query, fmt.Errorf("invalid value type for IN filter on field %s", filter.Field)
			}
		case "NOT IN":
			if vals, ok := filter.Value.([]interface{}); ok {
				query = query.Where(squirrel.NotEq{string(filter.Field): vals})
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
