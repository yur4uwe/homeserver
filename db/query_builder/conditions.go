package querybuilder

func (q *QueryBuilder) Where(clause string) *QueryBuilder {
	q.conditions = clause
	return q
}
