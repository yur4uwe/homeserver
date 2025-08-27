package auth_test

import (
	"context"
	"errors"
	"homeserver/db/repositories"
	"homeserver/pkg/handlers/web/auth"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutHandler(t *testing.T) {
	_, _, _, mockTokR, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, nil, nil, mockTokR, nil)
	mockErrResp := new(mocks.ErrorResponseMock)
	mockParser := new(mocks.MockParser)
	mockTemplate := new(mocks.MockTemplate)

	cleanupErrResp := helpers.MockFunction(&auth.ErrorResponse, mockErrResp.ErrorResponse, func() {})
	cleanupMust := helpers.MockFunction(&auth.Must, mockTemplate.Must, func() {})
	cleanupParseFiles := helpers.MockFunction(&auth.ParseFiles, mockTemplate.ParseFiles, func() {})
	cleanupParseTemplate := helpers.MockFunction(&auth.ParseTemplateToString, mockParser.ParseTemplateToString, func() {})
	defer cleanupErrResp()
	defer cleanupMust()
	defer cleanupParseFiles()
	defer cleanupParseTemplate()

	reqCont := middleware.DefaultContext
	reqCont.AuthToken = "token123"

	resetter := helpers.NewMocksResetter(
		&mockTokR.Mock,
		&mockErrResp.Mock,
		&mockParser.Mock,
		&mockTemplate.Mock,
	)

	handler := auth.LogoutHandler(mockDB)

	mockedTemplate, err := template.New("header").Parse("<header>{{.IsLoggedIn}}</header>")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		setupMocks     func()
		authToken      string
		wantCode       int
		wantBodySubstr string
	}{
		{
			name: "missing request context",
			setupMocks: func() {
				resetter()

				mockErrResp.On("ErrorResponse", mock.Anything, "Invalid request context", http.StatusInternalServerError).Return()
			},
			authToken:      "token123",
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "Invalid request context",
		},
		{
			name: "missing auth token",
			setupMocks: func() {
				resetter()
				mockErrResp.On("ErrorResponse", mock.Anything, "Unauthorized", http.StatusUnauthorized).Return()
			},
			authToken:      "",
			wantCode:       http.StatusUnauthorized,
			wantBodySubstr: "Unauthorized",
		},
		{
			name: "delete token error",
			setupMocks: func() {
				resetter()
				mockErrResp.On("ErrorResponse", mock.Anything, "Failed to delete token: db error", http.StatusInternalServerError).Return()
				mockTokR.On("DeleteToken", mock.Anything, "token123").Return(errors.New("db error"))
			},
			authToken:      "token123",
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "Failed to delete token",
		},
		{
			name: "header template error",
			setupMocks: func() {
				resetter()
				mockTokR.On("DeleteToken", mock.Anything, "token123").Return(nil)
				mockParser.On("ParseTemplateToString", mock.Anything, mock.Anything).Return("", errors.New("header error"))
				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError)
			},
			authToken:      "token123",
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "Internal Server Error: header error",
		},
		{
			name: "success logout",
			setupMocks: func() {
				resetter()
				mockTokR.On("DeleteToken", mock.Anything, "token123").Return(nil)
				mockParser.On("ParseTemplateToString", mock.Anything, mock.Anything).Return("<header>", nil)
				mockTemplate.On("Must", mock.Anything, mock.Anything).Return(mockedTemplate)
				mockTemplate.On("ParseFiles", mock.Anything).Return(mockedTemplate, nil)
			},
			authToken:      "token123",
			wantCode:       http.StatusOK,
			wantBodySubstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			if tt.name != "missing request context" {
				reqCont.AuthToken = tt.authToken
				req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, reqCont))
			}

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
			if tt.wantBodySubstr != "" {
				assert.Contains(t, rec.Body.String(), tt.wantBodySubstr)
			}
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, "/", rec.Header().Get("HX-Push-Url"))
				assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
			}
		})
	}
}
