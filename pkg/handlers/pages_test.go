package handlers

import (
	"context"
	"fmt"
	"homeserver/pkg/middleware"
	"homeserver/test"
	"homeserver/test/mocks"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWelcomePageHandler(t *testing.T) {
	mockMust := new(mocks.MockMust)
	cleanupMust := test.MockFunction(&Must, mockMust.Must, func() {})
	mockParseFiles := new(mocks.MockParseFiles)
	cleanupParseFiles := test.MockFunction(&ParseFiles, mockParseFiles.ParseFiles, func() {})
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})

	defer func() {
		cleanupMust()
		cleanupParseFiles()
		cleanupErrorResponse()
	}()

	mockTemplate, err := template.New("welcome").Parse("<html><body>Welcome</body></html>")
	assert.NoError(t, err)

	t.Run("Renders welcome page with valid context", func(t *testing.T) {
		mockMust.On("Must", mock.Anything).Return(mockTemplate)
		mockParseFiles.On("ParseFiles", mock.Anything).Return(mockTemplate, nil)

		modContext := middleware.DefaultContext
		modContext.AuthToken = "token"
		modContext.LogHistory = &mocks.MockLogHistory{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		WelcomePageHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Welcome")
	})

	t.Run("Returns error for missing context", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		WelcomePageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})
}

func TestUploadPageHandler(t *testing.T) {
	mockMust := new(mocks.MockMust)
	cleanupMust := test.MockFunction(&Must, mockMust.Must, func() {})
	mockParseFiles := new(mocks.MockParseFiles)
	cleanupParseFiles := test.MockFunction(&ParseFiles, mockParseFiles.ParseFiles, func() {})
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})

	defer func() {
		cleanupMust()
		cleanupParseFiles()
		cleanupErrorResponse()
	}()

	mockTemplate, err := template.New("upload").Parse("<html><body>Upload</body></html>")
	assert.NoError(t, err)

	t.Run("Renders upload page with valid context", func(t *testing.T) {
		mockMust.On("Must", mock.Anything).Return(mockTemplate)
		mockParseFiles.On("ParseFiles", mock.Anything).Return(mockTemplate, nil)

		modContext := middleware.DefaultContext
		modContext.LogHistory = &mocks.MockLogHistory{}

		req := httptest.NewRequest(http.MethodGet, "/upload", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		UploadPageHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Upload")
	})

	t.Run("Returns error for missing context", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/upload", nil)
		rec := httptest.NewRecorder()

		UploadPageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})
}

func TestStaticPageHandler(t *testing.T) {
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})

	defer cleanupErrorResponse()

	t.Run("Serves welcome page for root path", func(t *testing.T) {
		called := false
		cleanupWelcome := test.MockFunction(WelcomePageHandler, func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("WelcomePage"))
		}, func() {})
		defer cleanupWelcome()

		modContext := middleware.DefaultContext
		modContext.LogHistory = &mocks.MockLogHistory{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		StaticPageHandler(rec, req)

		assert.True(t, called, "WelcomePageHandler should be called")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "WelcomePage")
	})

	t.Run("Returns error for missing context", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		StaticPageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})

	t.Run("Returns not found for non-page path", func(t *testing.T) {
		isPageMock := new(IsPageMock)
		cleanupIsPage := test.MockFunction(&IsPage, isPageMock.IsPage, func() {})
		defer cleanupIsPage()

		modContext := middleware.DefaultContext
		modContext.LogHistory = &mocks.MockLogHistory{}

		req := httptest.NewRequest(http.MethodGet, "/notapage", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		StaticPageHandler(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Returns error for missing file", func(t *testing.T) {
		isPageMock := new(IsPageMock)
		cleanupIsPage := test.MockFunction(&IsPage, isPageMock.IsPage, func() {})
		defer cleanupIsPage()

		modContext := middleware.DefaultContext
		modContext.LogHistory = &mocks.MockLogHistory{}

		req := httptest.NewRequest(http.MethodGet, "/missingfile", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		// Mock os.ReadFile to return error
		cleanupReadFile := test.MockFunction(&ReadFile, func(filename string) ([]byte, error) {
			return nil, fmt.Errorf("file not found")
		}, func() {})
		defer cleanupReadFile()

		StaticPageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})

	// t.Run("Returns error for write failure", func(t *testing.T) {
	// 	cleanupIsPage := test.MockFunction(&IsPage, func(path string) bool { return true }, func() {})
	// 	defer cleanupIsPage()

	// 	modContext := middleware.DefaultContext
	// 	modContext.LogHistory = &mocks.MockLogHistory{}

	// 	req := httptest.NewRequest(http.MethodGet, "/writefail", nil)
	// 	rec := httptest.NewRecorder()
	// 	req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

	// 	// Mock os.ReadFile to return valid content
	// 	cleanupReadFile := test.MockFunction(&ReadFile, func(filename string) ([]byte, error) {
	// 		return []byte("content"), nil
	// 	}, func() {})
	// 	defer cleanupReadFile()

	// 	// Monkey patch rec.Write to simulate error
	// 	origWrite := rec.Write
	// 	rec.Write = func(p []byte) (int, error) {
	// 		return 0, fmt.Errorf("write error")
	// 	}

	// 	StaticPageHandler(rec, req)

	// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	// 	mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)

	// 	// Restore original Write method
	// 	rec.Write = origWrite
	// })
}
