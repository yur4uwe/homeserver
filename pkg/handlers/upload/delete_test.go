package upload

import (
	"errors"
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

func TestDeleteFile(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockErrRes := new(mocks.ErrorResponseMock)

	cleanupErrResp := test.MockFunction(&ErrorResponse, mockErrRes.ErrorResponse, func() {})
	defer cleanupErrResp()

	mockReqContext := middleware.DefaultContext
	mockReqContext.Requester = user_schema.
		NewUserBuilder().
		WithID(1).
		WithUsername("testuser")

	handler := DeleteFile(mockDB, mockReqContext)

	resetter := test.NewMocksResetter(&mockDB.Mock)

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

			mockDB.On("DeleteFile", mock.Anything, mock.Anything, mockReqContext.Requester.Username).Return(0, err)
			mockErrRes.On("ErrorResponse", mock.Anything, mock.Anything, tt.statusCode).Return()

			req := httptest.NewRequest(http.MethodDelete, "/delete?filename="+tt.filename, nil)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)

			if !tt.shouldErr {
				mockDB.AssertExpectations(t)
			}
		})
	}
}
