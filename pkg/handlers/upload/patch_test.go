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

func TestPatchFile(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockErrRes := new(mocks.ErrorResponseMock)

	cleanupErrResp := test.MockFunction(&ErrorResponse, mockErrRes.ErrorResponse, func() {})
	defer cleanupErrResp()

	mockReqContext := middleware.DefaultContext
	mockReqContext.Requester = user_schema.
		NewUserBuilder().
		WithID(1).
		WithUsername("testuser")

	handler := PatchFile(mockDB, mockReqContext)

	resetter := test.NewMocksResetter(&mockDB.Mock)

	table := []struct {
		name         string
		filename     string
		newContent   string
		statusCode   int
		shouldErr    bool
		expectedBody string
	}{
		{"valid file", "testfile.txt", "new+file.txt", http.StatusOK, false, "new file.txt"},
		{"missing file", "folder/missingfile.txt", "new+file.txt", http.StatusInternalServerError, true, "Error renaming file"},
		{"empty filename", "", "new+file.txt", http.StatusBadRequest, true, "Missing filename parameter"},
		{"empty new filename", "testfile.txt", "", http.StatusBadRequest, true, "Missing filename parameter"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resetter()

			var err error = nil
			if tt.shouldErr {
				err = errors.New("file not found")
			}

			mockDB.On("RenameFile",
				mock.Anything, mock.AnythingOfType("string"),
				mock.AnythingOfType("string"), mockReqContext.Requester,
			).Return("new file.txt", err)
			mockErrRes.On("ErrorResponse", mock.Anything, mock.Anything, tt.statusCode).Return()

			req := httptest.NewRequest(
				http.MethodPatch,
				"/patch?new_name="+tt.filename+"&old_name="+tt.newContent,
				nil,
			)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
