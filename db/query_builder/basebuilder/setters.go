package basebuilder

import "homeserver/constants/qbconsts"

func (b *BaseBuilder) From(table string) *BaseBuilder {
	b.table = table
	return b
}

func (b *BaseBuilder) Where(where string) *BaseBuilder {
	b.where = where
	return b
}

func (b *BaseBuilder) Returning(cols ...string) *BaseBuilder {
	for _, col := range cols {
		b.returning = append(b.returning, qbconsts.Column{Name: col})
	}
	return b
}

func (b *BaseBuilder) Raw(sql string) *BaseBuilder {
	b.raw = sql
	return b
}
