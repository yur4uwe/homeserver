package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		status         int
		expectedHeader string
	}{
		{
			name:           "Error response with status 400",
			message:        "Bad Request",
			status:         http.StatusBadRequest,
			expectedHeader: `{"raise-alert": {"title": "Error", "value": "Bad Request", "status": "error", "icon": "x"}}`,
		},
		{
			name:           "Error response with status 500",
			message:        "Internal Server Error",
			status:         http.StatusInternalServerError,
			expectedHeader: `{"raise-alert": {"title": "Error", "value": "Internal Server Error", "status": "error", "icon": "x"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ErrorResponse(w, tt.message, tt.status)

			assert.Equal(t, tt.status, w.Code, "Expected status code to match")
			assert.Equal(t, tt.expectedHeader, w.Header().Get("HX-Trigger"), "Expected HX-Trigger header to match")
			// Body should be empty
			assert.Equal(t, "", w.Body.String(), "Expected response body to be empty")
		})
	}
}
