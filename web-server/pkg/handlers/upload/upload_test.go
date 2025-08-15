package upload

import (
	"bytes"
	"context"
	"errors"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/test"
	"homeserver/test/mocks"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDownloadFile(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockErrResp := new(mocks.ErrorResponseMock)
	defer test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})()

	reqCont := middleware.DefaultContext
	reqCont.Requester = user_schema.NewUserBuilder().
		WithID(1).
		WithUsername("testuser")

	resetter := test.NewMocksResetter(&mockDB.Mock)

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
		mockDB.On("SaveFile", mock.Anything, mock.AnythingOfType("*file_schema.UserFile"), mock.Anything, reqCont.Requester).Return(1, "myfile.txt", nil)

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()

		handler := DownloadFile(mockDB, reqCont)
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

		handler := DownloadFile(mockDB, reqCont)
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error retrieving file")
	})

	t.Run("error saving file", func(t *testing.T) {
		resetter()
		body, contentType := createMultipart("myfile.txt", "")
		mockDB.On("SaveFile", mock.Anything, mock.AnythingOfType("*file_schema.UserFile"), mock.Anything, reqCont.Requester).Return(0, "", errors.New("save error"))

		mockErrResp.On("ErrorResponse", mock.Anything, "Error saving file", http.StatusInternalServerError).Return()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()

		handler := DownloadFile(mockDB, reqCont)
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

		handler := DownloadFile(mockDB, reqCont)
		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error retrieving file")
	})
}
