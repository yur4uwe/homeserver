package actions

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/pkg/handlers/web/actions"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRenameFileHandler(t *testing.T) {
	_, mockFR, _, _, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, mockFR, nil, nil, nil)

	mockErrResp := new(mocks.ErrorResponseMock)
	cleanupErrResp := helpers.MockFunction(&actions.ErrorResponse, mockErrResp.ErrorResponse, func() {})

	mockNewAlert := new(mocks.MockParser)
	cleanupNewAlert := helpers.MockFunction(&actions.NewAlert, mockNewAlert.NewAlert, func() {})

	resetter := helpers.NewMocksResetter(
		&mockFR.Mock,
		&mockErrResp.Mock,
		&mockNewAlert.Mock,
	)

	defer func() {
		cleanupErrResp()
		cleanupNewAlert()
	}()

	mockRequester := uschema.New().
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
			desiredName:  "new",
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
			desiredName:  "new",
			renameErr:    nil,
			renameResult: "new",
			wantCode:     http.StatusOK,
			wantBody:     "new.txt",
			relativePath: []string{"path/to/file"},
			is_folder:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetter()

			newf := &fschema.UserFile{
				Filename:  tt.renameResult,
				Extension: "txt",
			}

			mockFR.On(
				"UpdateFileEntry",
				context.Background(),
				mock.AnythingOfType("*fschema.UserFile"),
				mock.AnythingOfType("*fschema.UserFile"),
				mockRequester,
			).Return(newf, tt.renameErr)
			mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return(tt.wantBody)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/rename?desired_name=%s", tt.desiredName), nil)
			rec := httptest.NewRecorder()

			additionalContext := middleware.DefaultContext
			additionalContext.Requester = mockRequester

			actions.RenameFileHandler(rec, req, mockDB, additionalContext, tt.relativePath, tt.fileName, tt.is_folder)

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, rec.Body.String())
		})
	}
}
