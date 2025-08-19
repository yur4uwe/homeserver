package querybuilder

import (
	"homeserver/constants/qbconsts"
	"strconv"
	"strings"
)

func (q *QueryBuilder) buildInsert() string {
	if q.queryType != qbconsts.INSERT {
		panic("QueryBuilder is not in INSERT mode")
	}

	query := "INSERT INTO " + q.table + " ("
	query += strings.Join(q.columns, ", ") + ") VALUES ("
	for i := 0; i < q.valuesNum; i++ {
		if i > 0 {
			query += ", "
		}
		query += "$" + strconv.Itoa(i+1)
	}
	query += ")"
	return query
}

func (q *QueryBuilder) Insert(columns ...string) *QueryBuilder {
	q.queryType = qbconsts.INSERT
	q.columns = columns
	return q
}

func (q *QueryBuilder) Into(table string) *QueryBuilder {
	q.table = table
	return q
}

func (q *QueryBuilder) Values(valuesNum int) *QueryBuilder {
	q.valuesNum = valuesNum
	return q
}
