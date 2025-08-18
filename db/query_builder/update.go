package querybuilder

import (
	"homeserver/constants/qbconsts"
	"strconv"
	"strings"
)

func (q *QueryBuilder) buildUpdate() string {
	query := "UPDATE " + q.table + " SET "
	for i, col := range q.columns {
		if i > 0 {
			query += ", "
		}
		query += col + " = $" + strconv.Itoa(i+1)
	}
	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
	}
	return query
}

func (q *QueryBuilder) Update(columns ...string) *QueryBuilder {
	q.queryType = qbconsts.UPDATE
	q.columns = columns
	return q
}

func (q *QueryBuilder) Set(columns ...string) *QueryBuilder {
	q.columns = columns
	return q
}
