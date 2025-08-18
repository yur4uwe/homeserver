package querybuilder

import "homeserver/constants/qbconsts"

func (q *QueryBuilder) Where(clause string) *QueryBuilder {
	q.conditions = append(q.conditions, clause)
	return q
}

func (q *QueryBuilder) And() *QueryBuilder {
	q.conditions = append(q.conditions, string(qbconsts.AND))
	return q
}

func (q *QueryBuilder) Or() *QueryBuilder {
	q.conditions = append(q.conditions, string(qbconsts.OR))
	return q
}

func (q *QueryBuilder) Not() *QueryBuilder {
	q.conditions = append(q.conditions, string(qbconsts.NOT))
	return q
}
