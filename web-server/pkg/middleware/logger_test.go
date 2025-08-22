package middleware

import (
	"context"
	"homeserver/test/helpers/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockhandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
	if !ok || reqCont == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("Request received")
	reqCont.LogHistory.ErrLog("An error occurred", "main.go:42")

	w.WriteHeader(http.StatusOK)
})

func TestLoggerMiddlewareUnit(t *testing.T) {
	mockLogHistory := new(mocks.MockLogHistory)
	mockLogHistory.On("Log", "Request received").Return()
	mockLogHistory.On("ErrLog", "An error occurred", "main.go:42").Return()
	mockLogHistory.On("UnloadLogs", mock.Anything).Return()

	mockContext := &AdditionalRequestContext{
		LogHistory: mockLogHistory,
	}

	logger := Logger(mockhandler)

	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), ReqContKey, mockContext))

	w := httptest.NewRecorder()

	logger.ServeHTTP(w, req)

	mockLogHistory.AssertCalled(t, "Log", "Request received")
	mockLogHistory.AssertCalled(t, "ErrLog", "An error occurred", "main.go:42")
	mockLogHistory.AssertCalled(t, "UnloadLogs", mock.Anything)

	assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
}

func TestLoggerMiddlewareIntegration(t *testing.T) {
	logHistory := &LogHistory{}

	mockContext := &AdditionalRequestContext{
		LogHistory: logHistory,
	}

	logger := Logger(mockhandler)

	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), ReqContKey, mockContext))

	w := httptest.NewRecorder()

	logger.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")

	assert.Len(t, logHistory.Entries, 2, "Log history should contain three entries")
	assert.Contains(t, logHistory.Entries[0], "Request received", "First log entry should contain 'Request received'")
	assert.Contains(t, logHistory.Entries[1], "[Err at main.go:42]", "Second log entry should contain error information")
}
