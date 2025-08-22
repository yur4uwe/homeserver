package middleware_test

import (
	"context"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/templates"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockHandler(contentType string, statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}
}

func TestInterceptor(t *testing.T) {
	mockLogHistory := new(mocks.MockLogHistory)
	mockLogHistory.On("Log", mock.Anything).Return()

	mockTemplate := new(mocks.MockTemplate)
	mockTemplate.On("ParseFiles", mock.Anything, mock.Anything).
		Return(mockIndexTemplate, nil)
	mockTemplate.On("ParseFiles", mock.Anything).
		Return(mockedTemplate, nil)

	cleanupParseFiles := helpers.MockFunction(&middleware.ParseFiles, mockTemplate.ParseFiles, func() {})
	defer cleanupParseFiles()

	mockErrorResponse := new(mocks.ErrorResponseMock)
	mockErrorResponse.On("ErrorResponse", mock.Anything, "Internal Server Error", http.StatusInternalServerError).Return()

	errorResponseCleanup := helpers.MockFunction(&middleware.ErrorResponse, mockErrorResponse.ErrorResponse, func() {})
	defer errorResponseCleanup()

	mockContext := &middleware.AdditionalRequestContext{
		LogHistory:    mockLogHistory,
		ToRestorePage: false,
		RerenderBody:  false,
		PageContent:   &templates.LayoutPageData{},
	}

	t.Run("Handles non-HTML content type", func(t *testing.T) {
		handler := mockHandler("application/json", http.StatusOK, `{"message": "success"}`)

		interceptor := middleware.Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Content-Type should be application/json")
		assert.Equal(t, `{"message": "success"}`, w.Body.String(), "Response body should match")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] Response status code: 200")
	})

	t.Run("Handles HTML content type without flags", func(t *testing.T) {
		handler := mockHandler("text/html; charset=utf-8", http.StatusOK, "<html><body>Hello, World!</body></html>")

		interceptor := middleware.Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"), "Content-Type should be text/html")
		assert.Equal(t, "<html><body>Hello, World!</body></html>", w.Body.String(), "Response body should match")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] Response status code: 200")
	})

	t.Run("Handles ToRestorePage flag", func(t *testing.T) {
		mockContext.ToRestorePage = true

		handler := mockHandler("text/html; charset=utf-8", http.StatusOK, "<div>Restored Content</div>")

		interceptor := middleware.Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] ToRestorePage is true, restoring page")
	})

	t.Run("Handles RerenderBody flag", func(t *testing.T) {
		mockContext.ToRestorePage = false
		mockContext.RerenderBody = true

		handler := mockHandler("text/html; charset=utf-8", http.StatusOK, "<div>Content</div>")

		interceptor := middleware.Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		assert.Equal(t, "body", w.Header().Get("HX-Retarget"), "HX-Retarget header should be set to body")
		assert.Equal(t, "innerHTML", w.Header().Get("HX-Reswap"), "HX-Reswap header should be set to innerHTML")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] RerenderBody is true, rendering body")
	})

	t.Run("Handles missing request context", func(t *testing.T) {
		handler := mockHandler("text/html; charset=utf-8", http.StatusOK, "<div>Content</div>")

		interceptor := middleware.Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Response status code should be 500 Internal Server Error")
		assert.Contains(t, w.Body.String(), "Internal Server Error", "Response body should be 'Internal Server Error'")

		mockErrorResponse.AssertCalled(t, "ErrorResponse", w, "Internal Server Error", http.StatusInternalServerError)
	})

	t.Run("Handles error response", func(t *testing.T) {
		handler := mockHandler("text/html; charset=utf-8", http.StatusInternalServerError, "<div>Error occurred</div>")

		interceptor := middleware.Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockContext))

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Response status code should be 500 Internal Server Error")
	})

	mockLogHistory.AssertExpectations(t)
	mockTemplate.AssertExpectations(t)
	mockErrorResponse.AssertExpectations(t)
}
