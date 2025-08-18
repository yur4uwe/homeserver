package querybuilder

import (
	"homeserver/constants/qbconsts"
	"strings"
)

func (q *QueryBuilder) buildSelect() string {
	if q.queryType != qbconsts.SELECT {
		panic("QueryBuilder is not in SELECT mode")
	}

	query := "SELECT " + strings.Join(q.columns, ", ") + " FROM " + q.table
	if q.conditions != "" {
		query += " WHERE " + q.conditions
	}
	return query
}

func (q *QueryBuilder) Select(columns ...string) *QueryBuilder {
	q.queryType = qbconsts.SELECT
	q.columns = columns
	return q
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
	q.table = table
	return q
}

func (q *QueryBuilder) OrderBy(order qbconsts.QueryOrder) *QueryBuilder {
	q.orderBy = order
	return q
}

func (q *QueryBuilder) Limit(n int) *QueryBuilder {
	q.limit = n
	return q
}

func (q *QueryBuilder) Offset(n int) *QueryBuilder {
	q.offset = n
	return q
}
