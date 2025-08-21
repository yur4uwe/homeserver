package basebuilder

import "homeserver/constants/qbconsts"

type BaseBuilder struct {
	table     string
	where     string
	returning []qbconsts.Column
	raw       string
}

func New() *BaseBuilder {
	return &BaseBuilder{}
}
