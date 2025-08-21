package qbselect

import (
	"homeserver/constants/qbconsts"
	"homeserver/db/query_builder/basebuilder"
)

type SelectQueryBuilder struct {
	bb       basebuilder.BaseBuilder
	columns  []qbconsts.Column
	joins    []qbconsts.Join
	orderBy  string
	limit    int
	offset   int
	groupBy  string
	having   string
	distinct bool
}

func NewSelect() *SelectQueryBuilder {
	return &SelectQueryBuilder{
		bb:       *basebuilder.New(),
		columns:  []qbconsts.Column{},
		joins:    []qbconsts.Join{},
		orderBy:  "",
		limit:    0,
		offset:   0,
		groupBy:  "",
		having:   "",
		distinct: false,
	}
}
