package handlers_test

import (
	"context"
	"fmt"
	webhandlers "homeserver/pkg/handlers/web"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"homeserver/test/helpers/mocks/mocksos"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWelcomePageHandler(t *testing.T) {
	mockTemplates := new(mocks.MockTemplate)
	mockErrResp := new(mocks.ErrorResponseMock)

	cleanupMust := helpers.MockFunction(&webhandlers.Must, mockTemplates.Must, func() {})
	cleanupParseFiles := helpers.MockFunction(&webhandlers.ParseFiles, mockTemplates.ParseFiles, func() {})
	cleanupErrorResponse := helpers.MockFunction(&webhandlers.ErrorResponse, mockErrResp.ErrorResponse, func() {})

	defer func() {
		cleanupMust()
		cleanupParseFiles()
		cleanupErrorResponse()
	}()

	mockedTemplate, err := template.New("welcome").Parse("<html><body>Welcome</body></html>")
	assert.NoError(t, err)

	t.Run("Renders welcome page with valid context", func(t *testing.T) {
		mockTemplates.On("Must", mock.Anything, mock.Anything).Return(mockedTemplate)
		mockTemplates.On("ParseFiles", mock.Anything).Return(mockedTemplate, nil)

		modContext := middleware.DefaultContext
		modContext.AuthToken = "token"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		webhandlers.WelcomePageHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Welcome")
	})

	t.Run("Returns error for missing context", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		webhandlers.WelcomePageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})
}

func TestUploadPageHandler(t *testing.T) {
	mockTemplates := new(mocks.MockTemplate)
	mockErrResp := new(mocks.ErrorResponseMock)

	cleanupMust := helpers.MockFunction(&webhandlers.Must, mockTemplates.Must, func() {})
	cleanupParseFiles := helpers.MockFunction(&webhandlers.ParseFiles, mockTemplates.ParseFiles, func() {})
	cleanupErrorResponse := helpers.MockFunction(&webhandlers.ErrorResponse, mockErrResp.ErrorResponse, func() {})

	resetter := helpers.NewMocksResetter(
		&mockTemplates.Mock,
		&mockErrResp.Mock,
	)

	defer func() {
		cleanupMust()
		cleanupParseFiles()
		cleanupErrorResponse()
	}()

	mockedTemplate, err := template.New("upload").Parse("<html><body>Upload</body></html>")
	assert.NoError(t, err)

	t.Run("Renders upload page with valid context", func(t *testing.T) {
		resetter()

		mockTemplates.On("Must", mock.Anything, mock.Anything).Return(mockedTemplate)
		mockTemplates.On("ParseFiles", mock.Anything).Return(mockedTemplate, nil)

		modContext := middleware.DefaultContext

		req := httptest.NewRequest(http.MethodGet, "/upload", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		webhandlers.UploadPageHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Upload")
	})

	t.Run("Returns error for missing context", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/upload", nil)
		rec := httptest.NewRecorder()

		webhandlers.UploadPageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})
}

func TestStaticPageHandler(t *testing.T) {
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := helpers.MockFunction(&webhandlers.ErrorResponse, mockErrResp.ErrorResponse, func() {})

	mockOS := new(mocksos.MockOS)
	cleanupReadFile := helpers.MockFunction(&webhandlers.ReadFile, mockOS.ReadFile, func() {})

	isPageMock := new(IsPageMock)
	cleanupIsPage := helpers.MockFunction(&webhandlers.IsPage, isPageMock.IsPage, func() {})

	mockWelcomePage := new(WelcomeMock)
	cleanupWelcome := helpers.MockFunction(&webhandlers.WelcomePage, mockWelcomePage.WelcomePageHandler, func() {})

	resetter := helpers.NewMocksResetter(
		&mockErrResp.Mock,
		&mockOS.Mock,
		&isPageMock.Mock,
		&mockWelcomePage.Mock,
	)

	defer func() {
		cleanupErrorResponse()
		cleanupReadFile()
		cleanupIsPage()
		cleanupWelcome()
	}()

	t.Run("Serves static page with valid context", func(t *testing.T) {
		resetter()

		mockOS.On("ReadFile", mock.AnythingOfType("string")).Return([]byte("static content"), nil)
		isPageMock.On("IsPage", mock.Anything).Return(true)

		modContext := middleware.DefaultContext

		req := httptest.NewRequest(http.MethodGet, "/static", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		webhandlers.StaticPageHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "static content")
	})

	t.Run("Serves welcome page for root path", func(t *testing.T) {
		resetter()

		mockWelcomePage.On("WelcomePageHandler", mock.Anything, mock.Anything).Return()
		isPageMock.On("IsPage", mock.Anything).Return(true)

		modContext := middleware.DefaultContext

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		webhandlers.StaticPageHandler(rec, req)

		mockWelcomePage.AssertCalled(t, "WelcomePageHandler", rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "WelcomePage")
	})

	t.Run("Returns error for missing context", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		webhandlers.StaticPageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})

	t.Run("Returns not found for non-page path", func(t *testing.T) {
		resetter()

		isPageMock.On("IsPage", mock.Anything).Return(false)

		modContext := middleware.DefaultContext

		req := httptest.NewRequest(http.MethodGet, "/notapage", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		webhandlers.StaticPageHandler(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Returns error for missing file", func(t *testing.T) {
		resetter()

		mockOS.On("ReadFile", mock.Anything).Return(([]byte)(nil), fmt.Errorf("file not found"))
		isPageMock.On("IsPage", mock.Anything).Return(true)
		mockErrResp.On("ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError).Return()

		modContext := middleware.DefaultContext

		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rec := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		webhandlers.StaticPageHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, mock.AnythingOfType("string"), http.StatusInternalServerError)
	})
}
