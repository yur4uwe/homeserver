package querybuilder

import (
	"homeserver/constants/qbconsts"
	"strconv"
)

func (q *QueryBuilder) buildUpdate() string {
	if q.queryType != qbconsts.UPDATE {
		panic("QueryBuilder is not in UPDATE mode")
	}

	query := "UPDATE " + q.table + " SET "
	for i, col := range q.columns {
		if i > 0 {
			query += ", "
		}
		query += col + " = $" + strconv.Itoa(i+1)
	}
	if q.conditions != "" {
		query += " WHERE " + q.conditions
	}
	return query
}

func (q *QueryBuilder) Update(table string) *QueryBuilder {
	q.queryType = qbconsts.UPDATE
	q.table = table
	return q
}

func (q *QueryBuilder) Set(columns ...string) *QueryBuilder {
	q.columns = columns
	return q
}
