package qbupdate

import (
	"fmt"
	"homeserver/db/query_builder/basebuilder"
	"strings"
)

// Set adds a SET clause (e.g., Set("name = $1"))
func (qb *UpdateQueryBuilder) Set(setCols ...string) *UpdateQueryBuilder {
	init_len := len(qb.setColumns)
	for i, col := range setCols {
		qb.setColumns = append(qb.setColumns, fmt.Sprintf("%s = $%d", col, i+1+init_len))
	}
	return qb
}

func (qb *UpdateQueryBuilder) Returning(columns ...string) *UpdateQueryBuilder {
	qb.bb.Returning(columns...)
	return qb
}

func (qb *UpdateQueryBuilder) Where(whereClause string) *UpdateQueryBuilder {
	qb.bb.Where(whereClause)
	return qb
}

func (qb *UpdateQueryBuilder) From(table string) *UpdateQueryBuilder {
	qb.bb.From(table)
	return qb
}

func (qb *UpdateQueryBuilder) Raw(rawClause string) *UpdateQueryBuilder {
	qb.bb.Raw(rawClause)
	return qb
}

// Build assembles the UPDATE SQL query
func (qb *UpdateQueryBuilder) Build() string {
	if qb.bb.GetFrom() == "" {
		panic("No table specified for UPDATE")
	}
	if len(qb.setColumns) == 0 {
		panic("No SET columns specified for UPDATE")
	}

	query := fmt.Sprintf("UPDATE %s SET %s", qb.bb.GetFrom(), strings.Join(qb.setColumns, ", "))

	if where := qb.bb.GetWhere(); where != "" {
		query += " WHERE " + where
	}
	if returning := qb.bb.GetReturning(); len(returning) > 0 {
		query += " RETURNING " + basebuilder.ParseColumns(returning)
	}
	if raw := qb.bb.GetRaw(); raw != "" {
		query += " " + raw
	}
	return query
}

// New creates a new UpdateQueryBuilder
func New() *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		bb:         *basebuilder.New(),
		setColumns: []string{},
	}
}
