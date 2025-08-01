package utils

import (
	"homeserver/test_utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAlertGenerator is a mock implementation of the AlertGenerator
type MockAlertGenerator struct {
	mock.Mock
}

func (m *MockAlertGenerator) GenerateAlert(title, iconName, message, status string) string {
	args := m.Called(title, iconName, message, status)
	return args.String(0)
}

func TestErrorResponse(t *testing.T) {
	// Create a mock alert generator
	mockAlert := new(MockAlertGenerator)

	// Define the mock behavior
	mockAlert.On("GenerateAlert", "Error", "x", "Bad Request", AlertError).Return("Mock Alert: Bad Request")
	mockAlert.On("GenerateAlert", "Error", "x", "Internal Server Error", AlertError).Return("Mock Alert: Internal Server Error")

	cleanup := test_utils.MockFunction(&alertGenerator, mockAlert.GenerateAlert, func() {})
	defer cleanup()

	tests := []struct {
		name         string
		message      string
		status       int
		expectedBody string
	}{
		{
			name:         "Error response with status 400",
			message:      "Bad Request",
			status:       http.StatusBadRequest,
			expectedBody: "Mock Alert: Bad Request",
		},
		{
			name:         "Error response with status 500",
			message:      "Internal Server Error",
			status:       http.StatusInternalServerError,
			expectedBody: "Mock Alert: Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResponseRecorder to capture the response
			w := httptest.NewRecorder()

			// Call the ErrorResponse function
			ErrorResponse(w, tt.message, tt.status)

			// Assert the status code
			assert.Equal(t, tt.status, w.Code, "Expected status code to match")

			// Assert the response body
			assert.Contains(t, w.Body.String(), tt.expectedBody, "Expected response body to contain the mock alert")

			// Assert headers
			assert.Equal(t, "#modal-content", w.Header().Get("HX-Retarget"), "Expected HX-Retarget header to match")
			assert.Equal(t, "innerHTML", w.Header().Get("HX-Reswap"), "Expected HX-Reswap header to match")
			assert.Equal(t, "false", w.Header().Get("HX-Push-Url"), "Expected HX-Push-Url header to match")
			assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"), "Expected Content-Type header to match")
		})
	}

	// Assert that the mock was called as expected
	mockAlert.AssertExpectations(t)
}
