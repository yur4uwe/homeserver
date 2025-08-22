package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type ErrorResponseMock struct {
	mock.Mock
}

func (m *ErrorResponseMock) ErrorResponse(w http.ResponseWriter, message string, status int) {
	m.Called(w, message, status)
	w.WriteHeader(status)
	w.Write([]byte(message))
}
