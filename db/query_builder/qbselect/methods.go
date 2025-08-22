package qbselect

import (
	"homeserver/constants/qbconsts"
	"homeserver/db/query_builder/basebuilder"
	"strconv"
)

func (qb *SelectQueryBuilder) Columns(args ...string) *SelectQueryBuilder {

	for _, arg := range args {
		qb.columns = append(qb.columns, qbconsts.Column{Name: arg})
	}

	return qb
}

func (qb *SelectQueryBuilder) Col(args ...string) *SelectQueryBuilder {
	switch len(args) {
	case 1:
		qb.columns = append(qb.columns, qbconsts.Column{Name: args[0]})
	case 2:
		qb.columns = append(qb.columns, qbconsts.Column{Name: args[1], Table: args[0]})
	case 3:
		qb.columns = append(qb.columns, qbconsts.Column{Name: args[1], Table: args[0], Alias: args[2]})
	default:
		panic("Col method requires 1 to 3 arguments")
	}
	return qb
}

func (qb *SelectQueryBuilder) Ret(args ...string) *SelectQueryBuilder {
	qb.bb.Returning(args...)
	return qb
}

func (qb *SelectQueryBuilder) From(table string) *SelectQueryBuilder {
	qb.bb.From(table)
	return qb
}

func (qb *SelectQueryBuilder) Join(joinType, table, on string) *SelectQueryBuilder {
	qb.joins = append(qb.joins, qbconsts.Join{Type: joinType, Table: table, On: on})
	return qb
}

func (qb *SelectQueryBuilder) Where(where string) *SelectQueryBuilder {
	qb.bb.Where(where)
	return qb
}

func (qb *SelectQueryBuilder) OrderBy(order string) *SelectQueryBuilder {
	qb.orderBy = order
	return qb
}

func (qb *SelectQueryBuilder) Limit(n int) *SelectQueryBuilder {
	qb.limit = n
	return qb
}

func (qb *SelectQueryBuilder) Offset(n int) *SelectQueryBuilder {
	qb.offset = n
	return qb
}

func (qb *SelectQueryBuilder) GroupBy(group string) *SelectQueryBuilder {
	qb.groupBy = group
	return qb
}

func (qb *SelectQueryBuilder) Having(having string) *SelectQueryBuilder {
	qb.having = having
	return qb
}

func (qb *SelectQueryBuilder) Distinct() *SelectQueryBuilder {
	qb.distinct = true
	return qb
}

func (qb *SelectQueryBuilder) Raw(sql string) *SelectQueryBuilder {
	qb.bb.Raw(sql)
	return qb
}

func (qb *SelectQueryBuilder) Build() string {
	// Build columns
	columns := "*"
	if len(qb.columns) > 0 {
		columns = basebuilder.ParseColumns(qb.columns)
	}

	// Handle DISTINCT
	selectClause := "SELECT "
	if qb.distinct {
		selectClause += "DISTINCT "
	}
	selectClause += columns

	// Build FROM
	query := selectClause + " FROM " + qb.bb.GetFrom()

	// Build JOINS
	for _, j := range qb.joins {
		query += " " + j.Type + " JOIN " + j.Table + " ON " + j.On
	}

	// Build WHERE
	if qb.bb.GetWhere() != "" {
		query += " WHERE " + qb.bb.GetWhere()
	}

	// Build GROUP BY
	if qb.groupBy != "" {
		query += " GROUP BY " + qb.groupBy
	}

	// Build HAVING
	if qb.having != "" {
		query += " HAVING " + qb.having
	}

	// Build ORDER BY
	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	// Build LIMIT
	if qb.limit > 0 {
		query += " LIMIT " + strconv.Itoa(qb.limit)
	}

	// Build OFFSET
	if qb.offset > 0 {
		query += " OFFSET " + strconv.Itoa(qb.offset)
	}

	if len(qb.bb.GetReturning()) > 0 {
		query += " RETURNING " + basebuilder.ParseColumns(qb.bb.GetReturning())
	}

	// Append raw SQL if present
	if qb.bb.GetRaw() != "" {
		query += " " + qb.bb.GetRaw()
	}

	return query
}
