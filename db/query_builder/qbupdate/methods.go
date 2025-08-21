package qbupdate

import (
	"fmt"
	"homeserver/db/query_builder/basebuilder"
	"strings"
)

// Set adds a SET clause (e.g., Set("name = $1"))
func (qb *UpdateQueryBuilder) Set(setClause string) *UpdateQueryBuilder {
	qb.setClauses = append(qb.setClauses, setClause)
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
	if len(qb.setClauses) == 0 {
		panic("No SET clauses specified for UPDATE")
	}
	query := fmt.Sprintf("UPDATE %s SET %s", qb.bb.GetFrom(), strings.Join(qb.setClauses, ", "))

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
		setClauses: []string{},
	}
}
