package querybuilder

import (
	"homeserver/constants/dbconsts"
	"homeserver/constants/qbconsts"
	"slices"
	"strings"
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
	var table dbconsts.DBTable

	switch q.table {
	case dbconsts.DBUsersTable.Name:
		table = *dbconsts.DBUsersTable
	case dbconsts.DBFilesTable.Name:
		table = *dbconsts.DBFilesTable
	case dbconsts.DBFoldersTable.Name:
		table = *dbconsts.DBFoldersTable
	default:
		panic("unknown table")
	}

	for _, col := range q.columns {
		if !slices.Contains(table.Fields, col) {
			panic("unknown column: " + col)
		}
	}

	for _, col := range q.returning {
		if !slices.Contains(table.Fields, col) {
			panic("unknown returning column: " + col)
		}
	}

	if q.queryType == "" {
		panic("empty query type")
	} else if q.table == "" {
		panic("empty table")
	} else if len(q.columns) == 0 && q.queryType != qbconsts.DELETE {
		panic("empty columns")
	}

	if len(q.columns) > len(table.Fields) {
		panic("too many columns")
	} else if len(q.returning) > len(table.Fields) {
		panic("too many returning columns")
	}

	query := ""
	switch q.queryType {
	case qbconsts.SELECT:
		query += q.buildSelect()
	case qbconsts.INSERT:
		query += q.buildInsert()
	case qbconsts.UPDATE:
		query += q.buildUpdate()
	case qbconsts.DELETE:
		query += q.buildDelete()
	default:
		panic("unknown query type")
	}

	if q.queryType != qbconsts.SELECT && len(q.returning) > 0 {
		query += " RETURNING " + strings.Join(q.returning, ", ")
	}

	return query
}
