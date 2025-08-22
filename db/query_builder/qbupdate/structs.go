package qbupdate

import (
	"homeserver/db/query_builder/basebuilder"
)

type UpdateQueryBuilder struct {
	bb         basebuilder.BaseBuilder
	setColumns []string
}
