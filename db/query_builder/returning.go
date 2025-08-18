package querybuilder

func (q *QueryBuilder) Returning(columns ...string) *QueryBuilder {
	q.returning = columns
	return q
}
