package middleware

import (
	"context"
	"homeserver/pkg/templates"
	"homeserver/test_utils"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInterceptor(t *testing.T) {
	mockLogHistory := new(MockLogHistory)
	mockLogHistory.On("Log", mock.Anything).Return()
	mockLogHistory.On("ErrLog", mock.Anything, mock.Anything).Return()

	mockParseFiles := new(MockParseFiles)
	mockParseFiles.On("ParseFiles", mock.Anything).Return(template.Must(template.New("").Parse("<html><body>{{.Content}}</body></html>")), nil)

	cleanupParseFiles := test_utils.MockFunction(&ParseFiles, mockParseFiles.ParseFiles, func() {})
	defer cleanupParseFiles()

	mockParseAsString := new(mockParseAsString)
	mockParseAsString.On("ParseTemplateToString", mock.Anything, mock.Anything).Return("Mocked String", nil)

	parseAsStringCleanup := test_utils.MockFunction(&ParseTemplateToString, mockParseAsString.ParseTemplateToString, func() {})
	defer parseAsStringCleanup()

	mockErrorResponse := new(MockErrorResponse)
	mockErrorResponse.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return("something")

	errorResponseCleanup := test_utils.MockFunction(&ErrorResponse, mockErrorResponse.ErrorResponse, func() {})
	defer errorResponseCleanup()

	mockContext := &AdditionalRequestContext{
		LogHistory:    mockLogHistory,
		ToRestorePage: false,
		RerenderBody:  false,
		PageContent:   &templates.LayoutPageData{},
	}

	t.Run("Handles non-HTML content type", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		})

		interceptor := Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Content-Type should be application/json")
		assert.Equal(t, `{"message": "success"}`, w.Body.String(), "Response body should match")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] Response status code: 200")
	})

	t.Run("Handles HTML content type without flags", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body>Hello, World!</body></html>"))
		})

		interceptor := Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"), "Content-Type should be text/html")
		assert.Equal(t, "<html><body>Hello, World!</body></html>", w.Body.String(), "Response body should match")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] Response status code: 200")
	})

	t.Run("Handles ToRestorePage flag", func(t *testing.T) {
		mockContext.ToRestorePage = true

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<div>Content</div>"))
		})

		interceptor := Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] ToRestorePage is true, restoring page")
	})

	t.Run("Handles RerenderBody flag", func(t *testing.T) {
		mockContext.ToRestorePage = false
		mockContext.RerenderBody = true

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<div>Content</div>"))
		})

		interceptor := Interceptor(handler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), ReqContKey, mockContext))

		w := httptest.NewRecorder()

		interceptor.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		assert.Equal(t, "body", w.Header().Get("HX-Retarget"), "HX-Retarget header should be set to body")
		assert.Equal(t, "innerHTML", w.Header().Get("HX-Reswap"), "HX-Reswap header should be set to innerHTML")
		mockLogHistory.AssertCalled(t, "Log", "[Interceptor] RerenderBody is true, rendering body")
	})

	mockErrorResponse.AssertNotCalled(t, "ErrorResponse", mock.Anything, mock.Anything, mock.Anything)
	mockLogHistory.AssertExpectations(t)
	mockParseFiles.AssertExpectations(t)
	mockParseAsString.AssertExpectations(t)
}
