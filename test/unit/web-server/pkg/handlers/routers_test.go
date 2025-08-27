package handlers_test

import (
	"context"
	"fmt"
	"homeserver/constants"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/internal/templates"
	handlers "homeserver/pkg/handlers/web"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"

	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadRouter(t *testing.T) {
	mockUpload := new(UploadMock)
	cleanupDelete := helpers.MockFunction(&handlers.DeleteFile, mockUpload.DeleteFile, func() {})
	cleanupDownload := helpers.MockFunction(&handlers.DownloadFile, mockUpload.DownloadFile, func() {})
	cleanupPatch := helpers.MockFunction(&handlers.PatchFile, mockUpload.PatchFile, func() {})

	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := helpers.MockFunction(&handlers.ErrorResponse, mockErrResp.ErrorResponse, func() {})

	defer func() {
		cleanupDelete()
		cleanupDownload()
		cleanupPatch()
		cleanupErrorResponse()
	}()

	var mockUploadHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	mockDownloadhandler := mockUploadHandler
	mockPatchHandler := mockUploadHandler

	t.Run("Hands over handling to DownloadFile on POST", func(t *testing.T) {
		mockUpload.On("DownloadFile", mock.Anything, mock.Anything).Return(mockDownloadhandler)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, middleware.DefaultContext))

		handler := handlers.UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockErrResp.AssertNotCalled(t, "ErrorResponse", mock.Anything, mock.Anything, mock.Anything)
		mockUpload.AssertCalled(t, "DownloadFile", mock.Anything, mock.Anything)
	})

	t.Run("Hands over handling to DeleteFile on DELETE", func(t *testing.T) {
		mockUpload.On("DeleteFile", mock.Anything, mock.Anything).Return(mockUploadHandler)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, middleware.DefaultContext))

		handler := handlers.UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockUpload.AssertCalled(t, "DeleteFile", mock.Anything, mock.Anything)
	})

	t.Run("Hands over handling to PatchFile on PATCH", func(t *testing.T) {
		mockUpload.On("PatchFile", mock.Anything, mock.Anything).Return(mockPatchHandler)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, middleware.DefaultContext))

		handler := handlers.UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockUpload.AssertCalled(t, "PatchFile", mock.Anything, mock.Anything)
	})

	t.Run("Returns Method Not Allowed for unsupported methods", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, "Method not allowed", http.StatusMethodNotAllowed).Return()

		req := httptest.NewRequest(http.MethodPut, "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, middleware.DefaultContext))

		handler := handlers.UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "Expected status code 405 Method Not Allowed")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Method not allowed", http.StatusMethodNotAllowed)
	})

	t.Run("Returns Internal Server Error for invalid context", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, "Invalid request context", http.StatusInternalServerError).Return()

		req := httptest.NewRequest("POST", "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		handler := handlers.UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500 Internal Server Error")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Invalid request context", http.StatusInternalServerError)
	})
}

func TestFileActionsRouter(t *testing.T) {
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := helpers.MockFunction(&handlers.ErrorResponse, mockErrResp.ErrorResponse, func() {})

	mockParsePath := new(mocks.ParsePathMock)
	cleanupParsePath := helpers.MockFunction(&handlers.ParsePathWithUName, mockParsePath.ParsePathWithUName, func() {})

	mockActions := new(ActionsMock)
	cleanupDelete := helpers.MockFunction(&handlers.DeleteFileHandler, mockActions.DeleteFileHandler, func() {})
	cleanupRename := helpers.MockFunction(&handlers.RenameFileHandler, mockActions.RenameFileHandler, func() {})

	resetter := helpers.NewMocksResetter(
		&mockErrResp.Mock,
		&mockParsePath.Mock,
		&mockActions.Mock,
	)

	defer func() {
		cleanupParsePath()
		cleanupErrorResponse()
		cleanupDelete()
		cleanupRename()
	}()

	t.Run("Hands over handling to DeleteFileHandler on DELETE", func(t *testing.T) {
		resetter()

		mockActions.On("DeleteFileHandler", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("[]string"), mock.AnythingOfType("string")).Return()
		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user", "file.txt"}, nil)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/file-actions/user/file.txt", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.FileActionsRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockActions.AssertCalled(t, "DeleteFileHandler", rec, req, mock.Anything, mock.Anything, []string{"user"}, "file.txt")
	})

	t.Run("Hands over handling to RenameFileHandler on PATCH", func(t *testing.T) {
		resetter()

		mockActions.On("RenameFileHandler", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, false).Return()
		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user", "file.txt"}, nil)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/file-actions/user/file.txt?desired_name=newfile.txt", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.FileActionsRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockActions.AssertCalled(t, "RenameFileHandler", rec, req, mock.Anything, mock.Anything, []string{"user"}, "file.txt", "newfile.txt", false)
	})

	t.Run("Returns Method Not Allowed for unsupported methods", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, "Method PUT not allowed", http.StatusMethodNotAllowed).Return()
		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user", "file.txt"}, nil)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodPut, "/api/v1/file-actions/user/file.txt", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.FileActionsRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "Expected status code 405 Method Not Allowed")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Method PUT not allowed", http.StatusMethodNotAllowed)
	})

	t.Run("Returns Bad Request for invalid file path", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, "Invalid file path", http.StatusBadRequest).Return()

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/file-actions/", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.FileActionsRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code 400 Bad Request")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Invalid file path", http.StatusBadRequest)
	})

	t.Run("Returns Bad Request for ParsePathWithUName error", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, "Parse error", http.StatusBadRequest).Return()
		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return(([]string)(nil), fmt.Errorf("Parse error"))

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/file-actions/user/file.txt", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.FileActionsRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code 400 Bad Request")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Parse error", http.StatusBadRequest)
	})

	t.Run("Returns Internal Server Error for missing request context", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, "Invalid request context", http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/file-actions/user/file.txt", nil)
		rec := httptest.NewRecorder()

		handler := handlers.FileActionsRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500 Internal Server Error")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Invalid request context", http.StatusInternalServerError)
	})
}

func TestSearchRouter(t *testing.T) {
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := helpers.MockFunction(&handlers.ErrorResponse, mockErrResp.ErrorResponse, func() {})

	mockParsePath := new(mocks.ParsePathMock)
	cleanupParsePath := helpers.MockFunction(&handlers.ParsePathWithUName, mockParsePath.ParsePathWithUName, func() {})

	mockSearch := new(SearchMock)
	cleanupSearch := helpers.MockFunction(&handlers.FileSearchHandler, mockSearch.FileSearchHandler, func() {})
	cleanupBrowse := helpers.MockFunction(&handlers.BrowseHandler, mockSearch.BrowseHandler, func() {})

	mockTemplate := new(mocks.MockTemplate)
	cleanupMust := helpers.MockFunction(&handlers.Must, mockTemplate.Must, func() {})
	cleanupParseFiles := helpers.MockFunction(&handlers.ParseFiles, mockTemplate.ParseFiles, func() {})

	mockedTemplate, err := template.New("test").Parse("<html><body>Test</body></html>")
	assert.NoError(t, err, "Failed to create mock template")

	resetter := helpers.NewMocksResetter(
		&mockErrResp.Mock,
		&mockParsePath.Mock,
		&mockSearch.Mock,
		&mockTemplate.Mock,
	)

	defer func() {
		cleanupParsePath()
		cleanupErrorResponse()
		cleanupSearch()
		cleanupBrowse()
		cleanupMust()
		cleanupParseFiles()
	}()

	t.Run("Hands over handling to FileSearchHandler if query isn't empty", func(t *testing.T) {
		resetter()

		mockSearch.On(
			"FileSearchHandler",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		).Return(&templates.SearchPageData{}, nil)

		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user"}, nil)

		mockTemplate.On(
			"Must",
			mock.Anything, mock.Anything,
		).Return(mockedTemplate)

		mockTemplate.On(
			"ParseFiles",
			mock.AnythingOfType("[]string"),
		).Return(mockedTemplate, nil)

		defer mockSearch.AssertExpectations(t)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/search/?query=test", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
	})

	t.Run("Hands over handling to BrowseHandler if query is empty", func(t *testing.T) {
		resetter()

		mockSearch.On(
			"BrowseHandler",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		).Return(&templates.SearchPageData{}, nil)

		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user", "parent", "folder"}, nil)

		mockTemplate.On(
			"Must",
			mock.Anything, mock.Anything,
		).Return(mockedTemplate)

		mockTemplate.On(
			"ParseFiles",
			mock.AnythingOfType("[]string"),
		).Return(mockedTemplate, nil)

		defer mockSearch.AssertExpectations(t)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/search/parent/folder", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
	})

	t.Run("Executes full template under certain conditions", func(t *testing.T) {
		resetter()

		mockSearch.On(
			"FileSearchHandler",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		).Return(&templates.SearchPageData{}, nil)

		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user"}, nil)

		mockTemplate.On(
			"Must",
			mock.Anything, mock.Anything,
		).Return(mockedTemplate)

		mockTemplate.On(
			"ParseFiles",
			[]string{constants.Rootd + "static/templates/search_base.html", constants.Rootd + "static/templates/search_results.html"},
		).Return(mockedTemplate, nil)

		defer mockSearch.AssertExpectations(t)
		defer mockTemplate.AssertExpectations(t)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")
		modContext.ToRestorePage = true

		req := httptest.NewRequest(http.MethodGet, "/search/?query=test", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		assert.Contains(t, rec.Header().Get("HX-Push-Url"), "/search/", "Expected HX-Push-Url to be '/search/'")
	})

	t.Run("Returns Internal Server Error for ParsePathWithUName error", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, "Internal Server Error: Failed to parse directory", http.StatusInternalServerError).Return()
		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return(([]string)(nil), fmt.Errorf("Parse error"))

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/search/", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500 Internal Server Error")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Internal Server Error: Failed to parse directory", http.StatusInternalServerError)
	})

	t.Run("Returns Internal Server Error if no additional context is present error", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, "Internal Server Error: Missing request context", http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodGet, "/search/", nil)
		rec := httptest.NewRecorder()

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500 Internal Server Error")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Internal Server Error: Missing request context", http.StatusInternalServerError)
	})

	t.Run("Returns Internal Server Error for FileSearchHandler error", func(t *testing.T) {
		resetter()

		mockErrResp.On(
			"ErrorResponse",
			mock.Anything,
			mock.AnythingOfType("string"),
			http.StatusInternalServerError,
		).Return()

		mockSearch.On(
			"FileSearchHandler",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		).Return(&templates.SearchPageData{}, fmt.Errorf("db query error"))

		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user"}, nil)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/search/?query=test", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500 Internal Server Error")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
	})

	t.Run("Returns Internal Server Error for BrowseHandler error", func(t *testing.T) {
		resetter()

		mockErrResp.On(
			"ErrorResponse",
			mock.Anything,
			mock.AnythingOfType("string"),
			http.StatusInternalServerError,
		).Return()

		mockSearch.On(
			"BrowseHandler",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		).Return(&templates.SearchPageData{}, fmt.Errorf("db query error"))

		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return([]string{"user"}, nil)

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/search/", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500 Internal Server Error")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
	})

	t.Run("Returns Bad Request for QueryUnescape error", func(t *testing.T) {
		resetter()

		mockErrResp.On("ErrorResponse", mock.Anything, "Internal Server Error: Failed to unescape directory", http.StatusBadRequest).Return()

		modContext := middleware.DefaultContext
		modContext.Requester = uschema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		req.URL.Path = "/search/%"

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := handlers.SearchRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code 400 Bad Request")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Internal Server Error: Failed to unescape directory", http.StatusBadRequest)
	})
}
