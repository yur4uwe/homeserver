package handlers

import (
	"context"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/test"
	"homeserver/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadRouter(t *testing.T) {
	mockUpload := new(UploadMock)
	cleanupDelete := test.MockFunction(&DeleteFile, mockUpload.DeleteFile, func() {})
	cleanupDownload := test.MockFunction(&DownloadFile, mockUpload.DownloadFile, func() {})
	cleanupPatch := test.MockFunction(&PatchFile, mockUpload.PatchFile, func() {})

	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})

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

		handler := UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockErrResp.AssertNotCalled(t, "ErrorResponse", mock.Anything, mock.Anything, mock.Anything)
		mockUpload.AssertCalled(t, "DownloadFile", nil, mock.Anything)
	})

	t.Run("Hands over handling to DeleteFile on DELETE", func(t *testing.T) {
		mockUpload.On("DeleteFile", mock.Anything, mock.Anything).Return(mockUploadHandler)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, middleware.DefaultContext))

		handler := UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockUpload.AssertCalled(t, "DeleteFile", nil, mock.Anything)
	})

	t.Run("Hands over handling to PatchFile on PATCH", func(t *testing.T) {
		mockUpload.On("PatchFile", mock.Anything, mock.Anything).Return(mockPatchHandler)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, middleware.DefaultContext))

		handler := UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockUpload.AssertCalled(t, "PatchFile", nil, mock.Anything)
	})

	t.Run("Returns Method Not Allowed for unsupported methods", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, "Method not allowed", http.StatusMethodNotAllowed).Return()

		req := httptest.NewRequest(http.MethodPut, "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, middleware.DefaultContext))

		handler := UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "Expected status code 405 Method Not Allowed")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Method not allowed", http.StatusMethodNotAllowed)
	})

	t.Run("Returns Internal Server Error for invalid context", func(t *testing.T) {
		mockErrResp.On("ErrorResponse", mock.Anything, "Invalid request context", http.StatusInternalServerError).Return()

		req := httptest.NewRequest("POST", "/api/v1/upload", nil)
		rec := httptest.NewRecorder()

		handler := UploadRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500 Internal Server Error")
		mockErrResp.AssertCalled(t, "ErrorResponse", mock.Anything, "Invalid request context", http.StatusInternalServerError)
	})
}

func TestFileActionsRouter(t *testing.T) {
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrorResponse := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})

	mockParsePath := new(mocks.ParsePathMock)
	cleanupParsePath := test.MockFunction(&ParsePathWithUName, mockParsePath.ParsePathWithUName, func() {})

	mockActions := new(ActionsMock)
	cleanupDelete := test.MockFunction(&DeleteFileHandler, mockActions.DeleteFileHandler, func() {})
	cleanupRename := test.MockFunction(&RenameFileHandler, mockActions.RenameFileHandler, func() {})
	cleanupOpen := test.MockFunction(&OpenFileHandler, mockActions.OpenFileHandler, func() {})

	defer func() {
		cleanupParsePath()
		cleanupErrorResponse()
		cleanupDelete()
		cleanupRename()
		cleanupOpen()
	}()

	t.Run("Hands over handling to DeleteFileHandler on DELETE", func(t *testing.T) {
		mockActions.On("DeleteFileHandler", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
		mockParsePath.On(
			"ParsePathWithUName",
			mock.Anything, mock.Anything,
		).Return("file.txt", []string{"user"}, nil)

		modContext := middleware.DefaultContext
		modContext.Requester = user_schema.NewUserBuilder().WithUsername("user")

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/file-actions/user/file.txt", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, modContext))

		handler := FileActionsRouter(nil)
		handler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
		mockActions.AssertCalled(t, "DeleteFileHandler", rec, req, nil, mock.Anything, []string{"user"}, "file.txt", false)
	})

	t.Run("Hands over handling to RenameFileHandler on PATCH", func(t *testing.T) {
		t.Skip("RenameFileHandler tests not implemented yet")
	})

	t.Run("Hands over handling to OpenFileHandler on GET", func(t *testing.T) {
		t.Skip("OpenFileHandler tests not implemented yet")
	})
}

func TestSearchRouter(t *testing.T) {
	t.Skip("SearchRouter tests not implemented yet")
	// This function would contain tests for the SearchRouter function
}
