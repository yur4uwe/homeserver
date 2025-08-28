package qbinsert

import (
	"homeserver/db/query_builder/basebuilder"
)

type InsertQueryBuilder struct {
	bb              basebuilder.BaseBuilder
	columns         []string
	conflictColumns []string
	conflictAction  string
	valuesNum       int
}
