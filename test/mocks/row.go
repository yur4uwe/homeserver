package mocks

import (
	"reflect"

	"github.com/stretchr/testify/mock"
)

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}

func MockRowReturn(expected ...any) *MockRow {
	mockRow := &MockRow{}
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		destSlice := args.Get(0).([]any)
		RunScan(destSlice, expected)
	})
	return mockRow
}

func RunScan(destSlice []any, values []any) {
	for i, v := range values {
		if i >= len(destSlice) {
			break
		}
		ptrVal := reflect.ValueOf(destSlice[i])
		if ptrVal.Kind() == reflect.Ptr {
			if v == nil {
				ptrVal.Elem().Set(reflect.Zero(ptrVal.Elem().Type()))
				continue
			}
			val := reflect.ValueOf(v)
			if val.Type().AssignableTo(ptrVal.Elem().Type()) {
				ptrVal.Elem().Set(val)
			} else if val.Type().ConvertibleTo(ptrVal.Elem().Type()) {
				ptrVal.Elem().Set(val.Convert(ptrVal.Elem().Type()))
			}
		}
	}
}
