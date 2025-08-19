package upload

import (
	"errors"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test"
	"homeserver/test/mocks"
	"homeserver/web-server/pkg/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteFile(t *testing.T) {
	_, mockFR, _, _, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, mockFR, nil, nil, nil)
	mockErrRes := new(mocks.ErrorResponseMock)

	cleanupErrResp := test.MockFunction(&ErrorResponse, mockErrRes.ErrorResponse, func() {})
	defer cleanupErrResp()

	mockReqContext := middleware.DefaultContext
	mockReqContext.Requester = uschema.
		NewUserBuilder().
		WithID(1).
		WithUsername("testuser")

	handler := DeleteFile(mockDB, mockReqContext)

	resetter := test.NewMocksResetter(
		&mockFR.Mock,
	)

	table := []struct {
		name         string
		filename     string
		statusCode   int
		shouldErr    bool
		expectedBody string
	}{
		{"valid file", "testfile.txt", http.StatusOK, false, ""},
		{"missing file", "folder/missingfile.txt", http.StatusInternalServerError, true, "Error deleting file"},
		{"empty filename", "", http.StatusBadRequest, true, "Missing filename parameter"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resetter()

			var err error = nil
			if tt.shouldErr {
				err = errors.New("file not found")
			}

			mockFR.On(
				"DeleteFile",
				mock.Anything,
				mock.AnythingOfType("*fschema.UserFile"),
				mockReqContext.Requester,
			).Return(err)
			mockErrRes.On("ErrorResponse", mock.Anything, mock.Anything, tt.statusCode).Return()

			req := httptest.NewRequest(http.MethodDelete, "/delete?filename="+tt.filename, nil)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)

			if !tt.shouldErr {
				mockFR.AssertExpectations(t)
			}
		})
	}
}
