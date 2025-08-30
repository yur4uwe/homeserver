package actions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/pkg/handlers/web/actions"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteFileHandler(t *testing.T) {
	_, mockFR, mockFolR, _, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, mockFR, mockFolR, nil, nil)
	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrResp := helpers.MockFunction(&actions.ErrorResponse, mockErrResp.ErrorResponse, func() {})
	defer cleanupErrResp()

	mockNewAlert := new(mocks.MockParser)
	cleanupNewAlert := helpers.MockFunction(&actions.NewAlert, mockNewAlert.NewAlert, func() {})
	defer cleanupNewAlert()

	resetter := helpers.NewMocksResetter(
		&mockFR.Mock,
		&mockFolR.Mock,
		&mockErrResp.Mock,
		&mockNewAlert.Mock,
	)

	mockRequester := uschema.New().
		WithID(1).
		WithUsername("testuser")

	additionalContext := middleware.DefaultContext
	additionalContext.Requester = mockRequester

	relativePath := []string{"folder"}
	fileName := "test"
	fileExtension := "txt"
	filePath := "folder"
	fileObj := fschema.New().
		WithFilename(fileName).
		WithExtension(fileExtension).
		WithPath(filePath).
		WithOwnerID(mockRequester.ID)

	folderName := "myfolder"
	folderPath := "folder"
	folderObj := folschema.New().
		WithName(folderName).
		WithPath(folderPath).
		WithOwnerID(mockRequester.ID)

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
			wantBody:    "ALERT: Folder Deleted check Folder 'myfolder' has been deleted successfully.",
			wantTrigger: "delete-success-folder-43",
		},
		{
			name:        "delete folder error",
			returnID:    0,
			returnErr:   fmt.Errorf("fail"),
			isFolder:    true,
			wantCode:    http.StatusInternalServerError,
			wantBody:    "Error deleting folder: fail",
			wantTrigger: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetter()

			if tt.isFolder {
				mockFolR.On("DeleteFolder", mock.Anything, folderObj, mockRequester).Return(tt.returnID, tt.returnErr)
				if tt.returnErr != nil {
					mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return(tt.wantBody)
				} else {
					mockNewAlert.On("NewAlert", "Folder Deleted", "check", fmt.Sprintf("Folder '%s' has been deleted successfully.", folderName), mock.Anything).Return(tt.wantBody)
				}
				req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
				rec := httptest.NewRecorder()
				actions.DeleteFileHandler(rec, req, mockDB, additionalContext, relativePath, folderName)
				assert.Equal(t, tt.wantCode, rec.Code)
				assert.Equal(t, tt.wantBody, rec.Body.String())
				if tt.wantTrigger != "" {
					assert.Equal(t, tt.wantTrigger, rec.Header().Get("HX-Trigger"))
				}
			} else {
				mockFR.On("DeleteFile", mock.Anything, fileObj, mockRequester).Return(tt.returnID, tt.returnErr)
				if tt.returnErr != nil {
					mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return(tt.wantBody)
				} else {
					mockNewAlert.On("NewAlert", "File Deleted", "check", fmt.Sprintf("File '%s.%s' has been deleted successfully.", fileName, fileExtension), mock.Anything).Return(tt.wantBody)
				}
				req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
				rec := httptest.NewRecorder()
				actions.DeleteFileHandler(rec, req, mockDB, additionalContext, relativePath, fileName+"."+fileExtension)
				assert.Equal(t, tt.wantCode, rec.Code)
				assert.Equal(t, tt.wantBody, rec.Body.String())
				if tt.wantTrigger != "" {
					assert.Equal(t, tt.wantTrigger, rec.Header().Get("HX-Trigger"))
				}
			}
		})
	}
}
