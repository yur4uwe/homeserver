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

func (b *BaseBuilder) Returning(col ...string) *BaseBuilder {
	switch len(col) {
	case 1:
		b.returning = append(b.returning, qbconsts.Column{Name: col[0]})
	case 2:
		b.returning = append(b.returning, qbconsts.Column{Name: col[1], Table: col[0]})
	case 3:
		b.returning = append(b.returning, qbconsts.Column{Name: col[1], Table: col[0], Alias: col[2]})
	default:
		panic("Returning method requires 1 to 3 arguments")
	}
	return b
}

func (b *BaseBuilder) Raw(sql string) *BaseBuilder {
	b.raw = sql
	return b
}
