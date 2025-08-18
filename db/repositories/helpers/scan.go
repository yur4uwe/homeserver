package helpers

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
)

func ScanRows(rows pgx.Rows, dest any) error {
	if rows == nil || !rows.Next() {
		return pgx.ErrNoRows
	} else if dest == nil {
		return fmt.Errorf("destination cannot be nil")
	}
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()
	fieldMap := map[string]int{}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag != "" {
			fieldMap[tag] = i
		}
	}

	descs := rows.FieldDescriptions()
	scanArgs := make([]any, len(descs))
	for i, desc := range descs {
		if idx, ok := fieldMap[desc.Name]; ok {
			scanArgs[i] = v.Field(idx).Addr().Interface()
		} else {
			var dummy any
			scanArgs[i] = &dummy
		}
	}
	if rows.Next() {
		return rows.Scan(scanArgs...)
	}
	return pgx.ErrNoRows
}
