package helpers

import (
	"errors"
	"reflect"

	"github.com/jackc/pgx/v5"
)

const (
	NilRowErr     = "rows cannot be nil"
	NilDestErr    = "destination cannot be nil"
	NotPtrToSlice = "dest must be a pointer to a slice"
)

func ScanRows(rows pgx.Rows, dest any) error {
	if rows == nil {
		return errors.New(NilRowErr)
	} else if dest == nil {
		return errors.New(NilDestErr)
	}

	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return errors.New(NotPtrToSlice)
	}
	sliceVal := destVal.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return errors.New(NotPtrToSlice)
	}
	elemType := sliceVal.Type().Elem()

	descs := rows.FieldDescriptions()
	fieldMap := map[string]int{}
	structType := elemType
	if elemType.Kind() == reflect.Ptr {
		structType = elemType.Elem()
	}
	for i := 0; i < structType.NumField(); i++ {
		tag := structType.Field(i).Tag.Get("db")
		if tag != "" {
			fieldMap[tag] = i
		}
	}

	for rows.Next() {
		var elem reflect.Value
		if elemType.Kind() == reflect.Ptr {
			elem = reflect.New(structType)
		} else {
			elem = reflect.New(structType).Elem()
		}

		structVal := elem
		if elemType.Kind() == reflect.Ptr {
			structVal = elem.Elem()
		}

		scanArgs := make([]any, len(descs))
		for i, desc := range descs {
			if idx, ok := fieldMap[desc.Name]; ok {
				scanArgs[i] = structVal.Field(idx).Addr().Interface()
			} else {
				var dummy any
				scanArgs[i] = &dummy
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}

		sliceVal.Set(reflect.Append(sliceVal, elem))
	}

	rows.Close()

	return nil
}
