package actions

import (
	"context"
	"fmt"
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

func TestRenameFileHandler(t *testing.T) {
	mockDB := new(mocks.MockDB)

	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrResp := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})

	mockNewAlert := new(mocks.MockParser)
	cleanupNewAlert := test.MockFunction(&NewAlert, mockNewAlert.NewAlert, func() {})

	resetter := test.NewMocksResetter(
		&mockDB.Mock,
		&mockErrResp.Mock,
		&mockNewAlert.Mock,
	)

	defer func() {
		cleanupErrResp()
		cleanupNewAlert()
	}()

	mockRequester := user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithID(1)

	tests := []struct {
		name         string
		renameErr    error
		renameResult string
		wantCode     int
		wantBody     string
		relativePath []string
		fileName     string
		desiredName  string
		is_folder    bool
	}{
		{
			name:         "empty desiredName",
			wantCode:     http.StatusBadRequest,
			wantBody:     "Desired name parameter is required for rename action",
			relativePath: []string{"path/to/file"},
			fileName:     "file.txt",
			desiredName:  "",
			is_folder:    false,
		},
		{
			name:         "rename error",
			fileName:     "file.txt",
			desiredName:  "new.txt",
			renameErr:    fmt.Errorf("fail"),
			renameResult: "",
			wantCode:     http.StatusInternalServerError,
			wantBody:     "Error renaming file: fail",
			relativePath: []string{"path/to/file"},
			is_folder:    false,
		},
		{
			name:         "success",
			fileName:     "old.txt",
			desiredName:  "new.txt",
			renameErr:    nil,
			renameResult: "new.txt",
			wantCode:     http.StatusOK,
			wantBody:     "new.txt",
			relativePath: []string{"path/to/file"},
			is_folder:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetter()

			mockDB.On(
				"RenameFile",
				context.Background(),
				mock.AnythingOfType("*file_schema.UserFile"),
				tt.desiredName,
				mockRequester,
			).Return(tt.renameResult, tt.renameErr)
			mockNewAlert.On("NewAlert", "File Renamed", "check", fmt.Sprintf("File '%s' has been renamed to '%s'.", tt.fileName, tt.renameResult), mock.Anything).Return("ALERT: File Renamed check File 'old.txt' has been renamed to 'new.txt'. 1")
			mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return(tt.wantBody)

			req := httptest.NewRequest(http.MethodPatch, "/rename", nil)
			rec := httptest.NewRecorder()

			additionalContext := middleware.DefaultContext
			additionalContext.Requester = mockRequester

			RenameFileHandler(rec, req, mockDB, additionalContext, tt.relativePath, tt.fileName, tt.desiredName, tt.is_folder)

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, rec.Body.String())
		})
	}
}
