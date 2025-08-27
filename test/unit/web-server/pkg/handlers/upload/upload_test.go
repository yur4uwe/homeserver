package upload_test

import (
	"bytes"
	"context"
	"errors"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/pkg/handlers/web/upload"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"homeserver/test/helpers/mocks/mocksos"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDownloadFile(t *testing.T) {
	_, mockFR, _, _, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, mockFR, nil, nil, nil)
	mockErrResp := new(mocks.ErrorResponseMock)
	defer helpers.MockFunction(&upload.ErrorResponse, mockErrResp.ErrorResponse, func() {})()

	reqCont := middleware.DefaultContext
	reqCont.Requester = uschema.NewUserBuilder().
		WithID(1).
		WithUsername("testuser")

	mockOS := new(mocksos.MockOS)
	defer helpers.MockFunction(&upload.MkdirAll, mockOS.MkdirAll, func() {})()
	defer helpers.MockFunction(&upload.Create, mockOS.Create, func() {})()
	defer helpers.MockFunction(&upload.Getenv, mockOS.Getenv, func() {})()

	mockIO := new(mocks.MockIO)
	defer helpers.MockFunction(&upload.Copy, mockIO.Copy, func() {})()

	mockInsertedFile := fschema.NewUserFileBuilder().
		WithID(1).
		WithFullName("myfile.txt").
		WithPath("testuser/docs").
		WithExtension("txt").
		WithSize(int64(len("filecontent"))).
		WithContentType("text/plain").
		WithOwnerID(reqCont.Requester.ID)

	resetter := helpers.NewMocksResetter(&mockFR.Mock)

	createMultipart := func(filename, folder string) (*bytes.Buffer, string) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", filename)
		io.Copy(part, bytes.NewReader([]byte("filecontent")))
		if folder != "" {
			writer.WriteField("folder", folder)
		}
		writer.Close()
		return body, writer.FormDataContentType()
	}

	t.Run("success upload", func(t *testing.T) {
		resetter()
		body, contentType := createMultipart("myfile.txt", "docs/myfile.txt")
		mockFR.On("CreateFile", mock.Anything, mock.Anything, reqCont.Requester).Return(mockInsertedFile, nil)
		mockIO.On("Copy", mock.Anything, mock.Anything).Return(int64(len("filecontent")), nil)
		mockOS.On("MkdirAll", mock.Anything, os.ModePerm).Return(nil)
		mockOS.On("Create", mock.Anything).Return(new(os.File), nil)
		mockOS.On("Getenv", "STORAGE_PATH").Return("/storage/path/")

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()

		handler := upload.DownloadFile(mockDB, reqCont)
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "myfile.txt", rec.Body.String())
		assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
	})

	t.Run("error retrieving file", func(t *testing.T) {
		resetter()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		mockErrResp.On("ErrorResponse", mock.Anything, "Error retrieving file", http.StatusBadRequest).Return()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		handler := upload.DownloadFile(mockDB, reqCont)
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error retrieving file")
	})

	t.Run("error saving file", func(t *testing.T) {
		resetter()
		body, contentType := createMultipart("myfile.txt", "")
		mockFR.On("CreateFile", mock.Anything, mock.Anything, reqCont.Requester).Return(mockInsertedFile, errors.New("error saving file"))

		mockErrResp.On("ErrorResponse", mock.Anything, "Error saving file", http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()

		handler := upload.DownloadFile(mockDB, reqCont)
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error saving file")
	})

	t.Run("request cancellation before file retrieval", func(t *testing.T) {
		resetter()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		mockErrResp.On("ErrorResponse", mock.Anything, "Error retrieving file", http.StatusBadRequest).Return()

		ctx, cancel := context.WithCancel(context.Background())
		req := httptest.NewRequest(http.MethodPost, "/upload", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		handler := upload.DownloadFile(mockDB, reqCont)
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error retrieving file")
	})
}
