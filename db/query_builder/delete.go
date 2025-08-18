package querybuilder

import (
	"homeserver/constants/qbconsts"
	"strings"
)

func (q *QueryBuilder) Delete() *QueryBuilder {
	q.queryType = qbconsts.DELETE
	return q
}

func (q *QueryBuilder) buildDelete() string {
	query := "DELETE FROM " + q.table
	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
	}
	return query
}
