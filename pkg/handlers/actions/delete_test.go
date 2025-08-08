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

func TestDeleteFileHandler(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrResp := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})
	defer cleanupErrResp()

	mockNewAlert := new(MockNewAlert)
	cleanupNewAlert := test.MockFunction(&NewAlert, mockNewAlert.NewAlert, func() {})
	defer cleanupNewAlert()

	resetter := test.NewMocksResetter(
		&mockDB.Mock,
		&mockErrResp.Mock,
		&mockNewAlert.Mock,
	)

	mockRequester := user_schema.NewUserBuilder().
		WithID(1).
		WithUsername("testuser")

	additionalContext := middleware.DefaultContext
	additionalContext.Requester = mockRequester

	tests := []struct {
		name        string
		returnID    int
		returnErr   error
		isFolder    bool
		wantCode    int
		wantBody    string
		wantTrigger string
	}{
		{
			name:        "delete error",
			returnID:    0,
			returnErr:   fmt.Errorf("fail"),
			isFolder:    false,
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Error deleting file: fail",
			wantTrigger: "",
		},
		{
			name:        "delete success file",
			returnID:    42,
			returnErr:   nil,
			isFolder:    false,
			wantCode:    http.StatusOK,
			wantBody:    "ALERT: File Deleted check File 'test.txt' has been deleted successfully. 1",
			wantTrigger: "delete-success-file-42",
		},
		{
			name:        "delete success folder",
			returnID:    43,
			returnErr:   nil,
			isFolder:    true,
			wantCode:    http.StatusOK,
			wantBody:    "ALERT: File Deleted check File 'test.txt' has been deleted successfully. 1",
			wantTrigger: "delete-success-folder-43",
		},
	}

	relativePath := []string{"folder"}
	fileName := "test.txt"
	filePath := "folder/test.txt"
	fileObj := file_schema.NewUserFileBuilder().
		WithFilename(fileName).
		WithPath(filePath).
		WithOwnerID(mockRequester.ID)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetter()

			mockDB.On("DeleteFile", mock.Anything, fileObj, mockRequester.Username).Return(tt.returnID, tt.returnErr)
			if tt.returnErr != nil {
				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return(tt.wantBody)
			} else {
				mockNewAlert.On("NewAlert", "File Deleted", "check", fmt.Sprintf("File '%s' has been deleted successfully.", fileName), mock.Anything).Return(tt.wantBody)
			}

			req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
			rec := httptest.NewRecorder()

			DeleteFileHandler(rec, req, mockDB, additionalContext, relativePath, fileName, tt.isFolder)
			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, rec.Body.String())
			if tt.wantTrigger != "" {
				assert.Equal(t, tt.wantTrigger, rec.Header().Get("HX-Trigger"))
			}
		})
	}
}
