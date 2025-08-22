package qbinsert

import (
	"fmt"
	"homeserver/db/query_builder/basebuilder"
	"strings"
)

// Into sets the table name for the insert.
func (qb *InsertQueryBuilder) Into(table string) *InsertQueryBuilder {
	qb.bb.From(table)
	return qb
}

func (qb *InsertQueryBuilder) Raw(sql string) *InsertQueryBuilder {
	qb.bb.Raw(sql)
	return qb
}

// Columns sets the columns for the insert.
func (qb *InsertQueryBuilder) Columns(cols ...string) *InsertQueryBuilder {
	qb.columns = append(qb.columns, cols...)
	return qb
}

// Values adds a row of values for the insert.
func (qb *InsertQueryBuilder) ValuesNum(n int) *InsertQueryBuilder {
	qb.valuesNum = n
	return qb
}

func (qb *InsertQueryBuilder) Returning(cols ...string) *InsertQueryBuilder {
	qb.bb.Returning(cols...)
	return qb
}

// Build assembles the INSERT SQL query.
func (qb *InsertQueryBuilder) Build() string {
	table := qb.bb.GetFrom()
	if table == "" {
		panic("No table specified for INSERT")
	}
	if len(qb.columns) == 0 {
		panic("No columns specified for INSERT")
	}
	if qb.valuesNum == 0 {
		qb.valuesNum = 1
	}

	columns := "(" + strings.Join(qb.columns, ", ") + ")"
	valuePlaceholders := make([]string, qb.valuesNum)
	argCounter := 1
	for i := 0; i < qb.valuesNum; i++ {
		placeholders := make([]string, len(qb.columns))
		for j := range qb.columns {
			placeholders[j] = fmt.Sprintf("$%d", argCounter)
			argCounter++
		}
		valuePlaceholders[i] = "(" + strings.Join(placeholders, ", ") + ")"
	}

	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", table, columns, strings.Join(valuePlaceholders, ", "))

	if returning := qb.bb.GetReturning(); len(returning) > 0 {
		query += " RETURNING " + basebuilder.ParseColumns(returning)
	}
	if raw := qb.bb.GetRaw(); raw != "" {
		query += " " + raw
	}
	return query
}

// New creates a new InsertQueryBuilder.
func New() *InsertQueryBuilder {
	return &InsertQueryBuilder{
		bb: *basebuilder.New(),
	}
}
