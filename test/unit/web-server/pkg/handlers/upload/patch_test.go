package upload_test

import (
	"errors"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strings"

	"homeserver/pkg/handlers/web/upload"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils/generics"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPatchFile(t *testing.T) {
	_, mockFR, _, _, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, mockFR, nil, nil, nil)
	mockErrRes := new(mocks.ErrorResponseMock)

	cleanupErrResp := helpers.MockFunction(&upload.ErrorResponse, mockErrRes.ErrorResponse, func() {})
	defer cleanupErrResp()

	mockReqContext := middleware.DefaultContext
	mockReqContext.Requester = uschema.
		New().
		WithID(1).
		WithUsername("testuser")

	handler := upload.PatchFile(mockDB, mockReqContext)

	resetter := helpers.NewMocksResetter(&mockFR.Mock)

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

			name := strings.Split(generics.Must(url.QueryUnescape(tt.newContent)).(string), ".")

			newf := &fschema.UserFile{}

			if len(name) == 2 {
				newf.Filename = name[0]
				newf.Extension = name[1]
			} else if len(name) == 1 {
				newf.Filename = name[0]
			}

			mockFR.On("UpdateFileEntry",
				mock.Anything,
				mock.AnythingOfType("*fschema.UserFile"),
				mock.AnythingOfType("*fschema.UserFile"),
				mockReqContext.Requester,
			).Return(newf, err)
			mockErrRes.On("ErrorResponse", mock.Anything, mock.Anything, tt.statusCode).Return()

			req := httptest.NewRequest(
				http.MethodPatch,
				"/patch?new_name="+tt.newContent+"&old_name="+tt.filename,
				nil,
			)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
