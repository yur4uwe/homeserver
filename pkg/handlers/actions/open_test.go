package actions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/test"
	"homeserver/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHttp struct {
	mock.Mock
}

func (m *MockHttp) ServeFile(w http.ResponseWriter, r *http.Request, filePath string) {
	m.Called(w, r, filePath)
}

func TestOpenFileHandler(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrResp := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})
	defer cleanupErrResp()

	mockServeFile := new(MockHttp)
	cleanupServeFile := test.MockFunction(&ServeFile, mockServeFile.ServeFile, func() {})
	defer cleanupServeFile()

	resetter := test.NewMocksResetter(
		&mockDB.Mock,
		&mockErrResp.Mock,
	)

	mockRequester := user_schema.NewUserBuilder().
		WithID(1).
		WithUsername("testuser")

	additionalContext := middleware.DefaultContext
	additionalContext.Requester = mockRequester

	relativePath := []string{"folder"}
	fileName := "test.txt"
	filePath := "folder"
	fileObj := file_schema.NewUserFileBuilder().
		WithFilename("test").
		WithExtension("txt").
		WithPath(filePath).
		WithOwnerID(mockRequester.ID)

	t.Run("success", func(t *testing.T) {
		resetter()

		mockDB.On("GetFile", mock.Anything, mock.Anything, mockRequester).Return(fileObj, nil)
		mockServeFile.On("ServeFile", mock.Anything, mock.Anything, fileObj.Path).Return()

		req := httptest.NewRequest(http.MethodGet, "/open", nil)
		rec := httptest.NewRecorder()

		OpenFileHandler(rec, req, mockDB, additionalContext, relativePath, fileName)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error opening file", func(t *testing.T) {
		resetter()

		mockDB.On("GetFile", mock.Anything, mock.Anything, mockRequester).Return(fileObj, fmt.Errorf("fail"))
		mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return("Error opening file: fail")
		mockServeFile.On("ServeFile", mock.Anything, mock.Anything, fileObj.Path).Return()

		req := httptest.NewRequest(http.MethodGet, "/open", nil)
		rec := httptest.NewRecorder()

		OpenFileHandler(rec, req, mockDB, additionalContext, relativePath, fileName)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "Error opening file: fail", rec.Body.String())
	})
}
