package basebuilder

import (
	"homeserver/constants/qbconsts"
	"strings"
)

func ParseColumns(cols []qbconsts.Column) string {
	parts := make([]string, 0, len(cols))
	for _, col := range cols {
		part := ""
		if col.Table != "" {
			part += col.Table + "."
		}
		part += col.Name
		if col.Alias != "" {
			part += " AS " + col.Alias
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, ", ")
}
