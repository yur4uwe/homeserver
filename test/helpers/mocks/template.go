package mocks

import (
	"html/template"

	"github.com/stretchr/testify/mock"
)

type MockTemplate struct {
	mock.Mock
}

func (m *MockTemplate) Must(value *template.Template, err error) *template.Template {
	args := m.Called(value, err)
	if err != nil {
		panic(err)
	}
	return args.Get(0).(*template.Template)
}

func (m *MockTemplate) ParseFiles(filenames ...string) (*template.Template, error) {
	args := m.Called(filenames)
	return args.Get(0).(*template.Template), args.Error(1)
}
