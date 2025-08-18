package querybuilder

import (
	"homeserver/constants/qbconsts"
)

type QueryBuilder struct {
	queryType  qbconsts.QueryType
	table      string
	columns    []string
	conditions string
	valuesNum  int
	returning  []string
	orderBy    qbconsts.QueryOrder
	limit      int
	offset     int
}

func New() *QueryBuilder {
	return &QueryBuilder{}
}

func (q *QueryBuilder) Reset() *QueryBuilder {
	q.queryType = ""
	q.table = ""
	q.columns = nil
	q.conditions = ""
	q.valuesNum = 0
	q.returning = nil
	q.orderBy = ""
	q.limit = 0
	q.offset = 0
	return q
}

func (q *QueryBuilder) Build() string {
	if q.queryType == "" {
		panic("empty query type")
	} else if q.table == "" {
		panic("empty table")
	} else if (len(q.columns) == 0 || q.valuesNum == 0) && q.queryType != qbconsts.DELETE {
		panic("empty columns")
	}

	switch q.queryType {
	case qbconsts.SELECT:
		return q.buildSelect()
	case qbconsts.INSERT:
		return q.buildInsert()
	case qbconsts.UPDATE:
		return q.buildUpdate()
	case qbconsts.DELETE:
		return q.buildDelete()
	default:
		return ""
	}
}
