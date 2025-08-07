package middleware

import (
	"context"
	"homeserver/test"
	"homeserver/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMiddlewareSuite(t *testing.T) {
	mockMust := new(mocks.MockMust)
	mockMust.On("Must", mock.Anything, mock.Anything).Return("Mocked String", nil)

	mustCleanup := test.MockFunction(&Must, mockMust.Must, func() {})
	defer mustCleanup()

	mockParse := new(MockParseAsString)
	mockParse.On("ParseTemplateToString", "../static/fragments/footer.html", nil).Return("Mocked Footer", nil)

	parseCleanup := test.MockFunction(&ParseTemplateToString, mockParse.ParseTemplateToString, func() {})
	defer parseCleanup()

	mockErrorResponse := new(MockErrorResponse)
	mockErrorResponse.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return()

	errorResponseCleanup := test.MockFunction(&ErrorResponse, mockErrorResponse.ErrorResponse, func() {})
	defer errorResponseCleanup()

	t.Run("TestContextConstructor", func(t *testing.T) {
		ctx := contextConstructor()

		assert.NotNil(t, ctx, "Context should not be nil")
		assert.Equal(t, *DefaultContext, *ctx, "Context should match the default context")
	})

	t.Run("TestWithAdditionalRequestContext", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
			assert.True(t, ok, "Request context should be of type AdditionalRequestContext")
			assert.NotNil(t, reqCont, "Request context should not be nil")
			assert.Equal(t, *DefaultContext, *reqCont, "Request context should match the default context")
			w.WriteHeader(http.StatusOK)
		})

		wrappedHandler := WithAdditionalRequestContext(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
	})

	t.Run("TestToRestorePage", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
			assert.True(t, ok, "Request context should be of type AdditionalRequestContext")
			assert.NotNil(t, reqCont, "Request context should not be nil")
			assert.True(t, reqCont.ToRestorePage, "ToRestorePage should be set to true")
			w.WriteHeader(http.StatusOK)
		})

		wrappedHandler := ToRestorePage(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("HX-History-Restore-Request", "true")
		req = req.WithContext(context.WithValue(req.Context(), ReqContKey, contextConstructor()))

		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		assert.Equal(t, "HX-History-Restore-Request, HX-Request", w.Header().Get("Vary"), "Vary header should be set correctly")
	})

	t.Run("TestToRestorePage with no context", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
			assert.False(t, ok, "Request context should not be of type AdditionalRequestContext")
			assert.Nil(t, reqCont, "Request context should be nil")
			w.WriteHeader(http.StatusOK)
		})

		wrappedHandler := ToRestorePage(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("HX-History-Restore-Request", "true")

		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Response status code should be 500 Internal Server Error")
		assert.NotEqual(t, "HX-History-Restore-Request, HX-Request", w.Header().Get("Vary"), "Vary header shouldn't be set")
	})
}
