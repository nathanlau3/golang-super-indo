package common

import (
	"fmt"
	"strings"
)

type FieldMatch struct {
	Column string
	Value  string
}

type Filter struct {
	Search       string
	SearchFields []string
	Fields       []FieldMatch
	SortBy       string
	Order        string
	Page         int
	Limit        int
	SortAllowed  map[string]string
	DefaultSort  string
}

func (f *Filter) AddField(column, value string) {
	if value != "" {
		f.Fields = append(f.Fields, FieldMatch{Column: column, Value: value})
	}
}

func (f *Filter) LoadFields(columns, values []string, allowed map[string]string) {
	for i, col := range columns {
		if i >= len(values) || values[i] == "" {
			continue
		}
		dbCol, ok := allowed[col]
		if !ok {
			continue
		}
		f.Fields = append(f.Fields, FieldMatch{Column: dbCol, Value: values[i]})
	}
}

func (f *Filter) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *Filter) BuildWhereClause(argIdx int) (string, []interface{}, int) {
	var conditions []string
	var args []interface{}

	if f.Search != "" && len(f.SearchFields) > 0 {
		var searchConds []string
		for _, field := range f.SearchFields {
			searchConds = append(searchConds, fmt.Sprintf("LOWER(COALESCE(%s, '')) LIKE LOWER($%d)", field, argIdx))
		}
		conditions = append(conditions, "("+strings.Join(searchConds, " OR ")+")")
		args = append(args, "%"+f.Search+"%")
		argIdx++
	}

	for _, field := range f.Fields {
		conditions = append(conditions, fmt.Sprintf("LOWER(COALESCE(%s, '')) = LOWER($%d)", field.Column, argIdx))
		args = append(args, field.Value)
		argIdx++
	}

	clause := ""
	if len(conditions) > 0 {
		clause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return clause, args, argIdx
}

func (f *Filter) BuildOrderClause() string {
	column := f.DefaultSort
	if column == "" {
		column = "created_at"
	}

	if f.SortBy != "" && f.SortAllowed != nil {
		if mapped, ok := f.SortAllowed[f.SortBy]; ok {
			column = mapped
		}
	}

	dir := "DESC"
	if f.Order == "asc" {
		dir = "ASC"
	}

	return fmt.Sprintf("ORDER BY %s %s", column, dir)
}

func (f *Filter) CacheKey(prefix string) string {
	parts := []string{prefix, f.Search}
	for _, field := range f.Fields {
		parts = append(parts, field.Column+":"+field.Value)
	}
	parts = append(parts, f.SortBy, f.Order, fmt.Sprintf("%d", f.Page), fmt.Sprintf("%d", f.Limit))
	return strings.Join(parts, ":")
}
