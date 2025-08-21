package qbdelete

import (
	"fmt"
	"homeserver/db/query_builder/basebuilder"
)

// From sets the table to delete from.
func (qb *DeleteQueryBuilder) From(table string) *DeleteQueryBuilder {
	qb.bb.From(table)
	return qb
}

// Where sets the WHERE clause.
func (qb *DeleteQueryBuilder) Where(where string) *DeleteQueryBuilder {
	qb.bb.Where(where)
	return qb
}

// Returning sets the RETURNING clause.
func (qb *DeleteQueryBuilder) Returning(cols ...string) *DeleteQueryBuilder {
	qb.bb.Returning(cols...)
	return qb
}

// Raw appends raw SQL at the end.
func (qb *DeleteQueryBuilder) Raw(sql string) *DeleteQueryBuilder {
	qb.bb.Raw(sql)
	return qb
}

// Build assembles the DELETE SQL query.
func (qb *DeleteQueryBuilder) Build() string {
	table := qb.bb.GetFrom()
	if table == "" {
		panic("No table specified for DELETE")
	}
	query := fmt.Sprintf("DELETE FROM %s", table)

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

// New creates a new DeleteQueryBuilder.
func New() *DeleteQueryBuilder {
	return &DeleteQueryBuilder{}
}
