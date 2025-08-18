package querybuilder

import (
	"homeserver/constants/qbconsts"
)

func (q *QueryBuilder) Delete() *QueryBuilder {
	q.queryType = qbconsts.DELETE
	return q
}

func (q *QueryBuilder) buildDelete() string {
	if q.queryType != qbconsts.DELETE {
		panic("QueryBuilder is not in DELETE mode")
	}

	query := "DELETE FROM " + q.table
	if q.conditions != "" {
		query += " WHERE " + q.conditions
	}
	return query
}
