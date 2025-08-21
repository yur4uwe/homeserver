package basebuilder

import "homeserver/constants/qbconsts"

func (bb *BaseBuilder) GetFrom() string {
	return bb.table
}

func (bb *BaseBuilder) GetWhere() string {
	return bb.where
}

func (bb *BaseBuilder) GetReturning() []qbconsts.Column {
	return bb.returning
}

func (bb *BaseBuilder) GetRaw() string {
	return bb.raw
}
