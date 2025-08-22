package actions

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"homeserver/web-server/pkg/middleware"
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
	cleanupErrResp := helpers.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})

	mockNewAlert := new(mocks.MockParser)
	cleanupNewAlert := helpers.MockFunction(&NewAlert, mockNewAlert.NewAlert, func() {})

	resetter := helpers.NewMocksResetter(
		&mockFR.Mock,
		&mockErrResp.Mock,
		&mockNewAlert.Mock,
	)

	defer func() {
		cleanupErrResp()
		cleanupNewAlert()
	}()

	mockRequester := uschema.NewUserBuilder().
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

			mockFR.On(
				"UpdateFileEntry",
				context.Background(),
				mock.AnythingOfType("*fschema.UserFile"),
				mock.AnythingOfType("map[string]interface {}"),
				mockRequester,
			).Return(tt.renameErr).Run(func(args mock.Arguments) {
				dest := args.Get(1).(*fschema.UserFile)
				if tt.renameResult != "" {
					dest.Filename = tt.renameResult
				}
			})
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
