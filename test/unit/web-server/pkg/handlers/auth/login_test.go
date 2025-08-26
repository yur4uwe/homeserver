package auth_test

import (
	"bytes"
	"context"
	"errors"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"homeserver/web-server/pkg/handlers/auth"
	"homeserver/web-server/pkg/middleware"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginHandler(t *testing.T) {
	mockUR, _, _, mockTokR, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(mockUR, nil, nil, mockTokR, nil)
	mockErrResp := new(mocks.ErrorResponseMock)
	mockLog := new(mocks.MockLogHistory)
	mockParser := new(mocks.MockParser)
	mockTemplate := new(mocks.MockTemplate)
	mockTokens := new(mocks.MockTokens)

	cleanupErrResp := helpers.MockFunction(&auth.ErrorResponse, mockErrResp.ErrorResponse, func() {})
	cleanupMust := helpers.MockFunction(&auth.Must, mockTemplate.Must, func() {})
	cleanupParseFiles := helpers.MockFunction(&auth.ParseFiles, mockTemplate.ParseFiles, func() {})
	cleanupParseTemplate := helpers.MockFunction(&auth.ParseTemplateToString, mockParser.ParseTemplateToString, func() {})
	cleanupNewToken := helpers.MockFunction(&auth.NewToken, mockTokens.NewToken, func() {})
	defer cleanupErrResp()
	defer cleanupMust()
	defer cleanupParseFiles()
	defer cleanupParseTemplate()
	defer cleanupNewToken()

	reqCont := middleware.DefaultContext

	resetter := helpers.NewMocksResetter(
		&mockUR.Mock,
		&mockTokR.Mock,
		&mockErrResp.Mock,
		&mockLog.Mock,
		&mockParser.Mock,
		&mockTemplate.Mock,
		&mockTokens.Mock,
	)

	handler := auth.LoginHandler(mockDB)

	mockedTemplate := template.Must(template.New("welcome").Parse("<h1>Welcome</h1>"))

	tests := []struct {
		name           string
		setupMocks     func()
		requestBody    string
		contentType    string
		context        bool
		wantCode       int
		wantBodySubstr string
		wantPushUrl    string
		wantReswap     string
		wantCookie     bool
	}{
		{
			name: "missing request context",
			setupMocks: func() {
				resetter()
				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()

			},
			requestBody:    "",
			contentType:    "",
			context:        false,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "Internal Server Error: Missing request context",
		},
		{
			name: "invalid login",
			setupMocks: func() {
				resetter()
			},
			requestBody:    "identificator=&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusOK,
			wantBodySubstr: "Identificator is required",
			wantPushUrl:    "",
			wantReswap:     "outerHTML",
		},
		{
			name: "user not found",
			setupMocks: func() {
				resetter()
				u := (*uschema.User)(nil)
				mockUR.On("Query", mock.Anything, mock.AnythingOfType("*uschema.User")).Return(u, errors.New("no rows in result set"))
			},
			requestBody:    "identificator=u&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusUnauthorized,
			wantBodySubstr: "Unauthorized: Invalid credentials",
		},
		{
			name: "db error",
			setupMocks: func() {
				resetter()

				mockUR.On("Query", mock.Anything, mock.AnythingOfType("*uschema.User")).Return((*uschema.User)(nil), errors.New("db error"))
				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
			},
			requestBody:    "identificator=u&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "",
		},
		{
			name: "email not verified",
			setupMocks: func() {
				resetter()

				user := &uschema.User{ID: 1, Username: "u", Email: "e", EmailVerified: false}
				mockUR.On("Query", mock.Anything, mock.AnythingOfType("*uschema.User")).Return(user, nil)
			},
			requestBody:    "identificator=u&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusForbidden,
			wantBodySubstr: "Forbidden: Email not verified",
		},
		{
			name: "token generation error",
			setupMocks: func() {
				resetter()

				user := &uschema.User{ID: 1, Username: "u", Email: "e", EmailVerified: true}
				mockUR.On("Query", mock.Anything, mock.AnythingOfType("*uschema.User")).Return(user, nil)
				mockTokens.On("NewToken").Return("", errors.New("token error"))
				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
			},
			requestBody:    "identificator=u&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "",
		},
		{
			name: "db create token error",
			setupMocks: func() {
				resetter()

				user := &uschema.User{ID: 1, Username: "u", Email: "e", EmailVerified: true}
				mockUR.On("Query", mock.Anything, mock.AnythingOfType("*uschema.User")).Return(user, nil)
				mockTokens.On("NewToken").Return("token123", nil)
				mockTokR.On("CreateToken", mock.Anything, 1, "token123").Return(errors.New("token db error"))
				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
			},
			requestBody:    "identificator=u&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "",
		},
		{
			name: "header template error",
			setupMocks: func() {
				resetter()

				user := &uschema.User{ID: 1, Username: "u", Email: "e", EmailVerified: true}
				mockUR.On("Query", mock.Anything, mock.AnythingOfType("*uschema.User")).Return(user, nil)
				mockTokens.On("NewToken").Return("token123", nil)
				mockTokR.On("CreateToken", mock.Anything, 1, "token123").Return(nil)
				mockParser.On("ParseTemplateToString", mock.Anything, mock.Anything).Return("", errors.New("header error"))
				mockLog.On("ErrLog", mock.Anything, mock.Anything).Return()
				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
			},
			requestBody:    "identificator=u&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "",
		},
		{
			name: "success login",
			setupMocks: func() {
				resetter()

				user := &uschema.User{ID: 1, Username: "u", Email: "e", EmailVerified: true}
				mockUR.On("Query", mock.Anything, mock.AnythingOfType("*uschema.User")).Return(user, nil)
				mockTokens.On("NewToken").Return("token123", nil)
				mockTokR.On("CreateToken", mock.Anything, 1, "token123").Return(nil)
				mockParser.On("ParseTemplateToString", mock.Anything, mock.Anything).Return("<header>", nil)
				mockTemplate.On("Must", mock.Anything, mock.Anything).Return(mockedTemplate)
				mockTemplate.On("ParseFiles", mock.Anything).Return(mockedTemplate, nil)
				mockTemplate.On("Execute", mock.Anything, mock.Anything).Return(nil)
			},
			requestBody:    "identificator=u&password=p",
			contentType:    "application/x-www-form-urlencoded",
			context:        true,
			wantCode:       http.StatusOK,
			wantPushUrl:    "/",
			wantBodySubstr: "",
			wantCookie:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			var req *http.Request
			if tt.contentType != "" {
				req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.requestBody))
				req.Header.Set("Content-Type", tt.contentType)
			} else {
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.requestBody))
			}
			if tt.context {
				req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, reqCont))
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantCode, rec.Code)
			if tt.wantBodySubstr != "" {
				assert.Contains(t, rec.Body.String(), tt.wantBodySubstr)
			}
			if tt.wantPushUrl != "" {
				assert.Equal(t, tt.wantPushUrl, rec.Header().Get("HX-Push-Url"))
			}
			if tt.wantReswap != "" {
				assert.Equal(t, tt.wantReswap, rec.Header().Get("HX-Reswap"))
			}
			if tt.wantCookie {
				cookies := rec.Result().Cookies()
				var found bool
				for _, c := range cookies {
					if c.Name == "authTokenHS" && c.Value == "token123" && c.Path == "/" && c.HttpOnly {
						found = true
						break
					}
				}
				assert.True(t, found, "authTokenHS cookie should be set")
			}
		})
	}
}
