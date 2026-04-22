package auth_test

import (
	"context"
	"errors"
	"homeserver/db/repositories"
	"homeserver/pkg/handlers/web/auth"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVerificationEmail struct {
	mock.Mock
}

func (m *MockVerificationEmail) SendVerificationEmail(form formatters.RegisterForm, token string) {
	m.Called(form, token)
}

func TestRegisterHandler(t *testing.T) {
	mockUR, _, _, mockTokR, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(mockUR, nil, nil, mockTokR, nil)
	mockErrResp := new(mocks.ErrorResponseMock)
	mockLog := new(mocks.MockLogHistory)
	mockTemplate := new(mocks.MockTemplate)
	mockVerificationEmail := new(MockVerificationEmail)
	mockTokens := new(mocks.MockTokens)
	mockParser := new(mocks.MockParser)

	cleanupNewToken := helpers.MockFunction(&auth.NewToken, mockTokens.NewToken, func() {})
	cleanupSendEmail := helpers.MockFunction(&auth.VerSuccess, mockVerificationEmail.SendVerificationEmail, func() {})
	cleanupErrResp := helpers.MockFunction(&auth.ErrorResponse, mockErrResp.ErrorResponse, func() {})
	cleanupMust := helpers.MockFunction(&auth.Must, mockTemplate.Must, func() {})
	cleanupParseFiles := helpers.MockFunction(&auth.ParseFiles, mockTemplate.ParseFiles, func() {})
	cleanupNewAlert := helpers.MockFunction(&auth.NewAlert, mockParser.NewAlert, func() {})
	defer cleanupErrResp()
	defer cleanupMust()
	defer cleanupParseFiles()
	defer cleanupSendEmail()
	defer cleanupNewToken()
	defer cleanupNewAlert()

	reqCont := middleware.DefaultContext

	resetter := helpers.NewMocksResetter(
		&mockUR.Mock,
		&mockTokR.Mock,
		&mockErrResp.Mock,
		&mockLog.Mock,
		&mockTemplate.Mock,
		&mockVerificationEmail.Mock,
		&mockTokens.Mock,
		&mockParser.Mock,
	)

	handler := auth.RegisterHandler(mockDB)

	tests := []struct {
		name           string
		setupMocks     func()
		requestBody    string
		context        bool
		wantCode       int
		wantBodySubstr string
		wantPushUrl    string
		wantReswap     string
	}{
		{
			name: "missing request context",
			setupMocks: func() {
				resetter()

				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
			},
			requestBody:    "",
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "Missing request context",
		},
		{
			name: "invalid registration",
			setupMocks: func() {
				resetter()
			},
			requestBody: "username=&email=e&password=p&confirm-password=p",

			context:        true,
			wantCode:       http.StatusOK,
			wantBodySubstr: "Username is required",
			wantPushUrl:    "",
			wantReswap:     "outerHTML",
		},
		{
			name: "db create user error",
			setupMocks: func() {
				resetter()

				form := formatters.RegisterForm{
					Username: "username",
					Email:    "e@example.com",
					Password: "p",
				}
				user := uschema.New().WithUsername(form.Username).WithEmail(form.Email).WithPassword(form.Password)

				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
				mockUR.On("Create", mock.Anything, user).Return((*uschema.User)(nil), errors.New("db error"))
			},
			requestBody: "username=username&email=e@example.com&password=p&confirm-password=p",

			context:        true,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "",
		},
		{
			name: "token generation error",
			setupMocks: func() {
				resetter()

				form := formatters.RegisterForm{
					Username: "u",
					Email:    "e@example.com",
					Password: "p",
				}
				user := uschema.New().WithUsername(form.Username).WithEmail(form.Email).WithPassword(form.Password)

				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
				mockUR.On("Create", mock.Anything, user).Return(&uschema.User{ID: 1, Username: "u", Email: "e"}, nil)
				mockTokens.On("NewToken", mock.Anything, mock.Anything).Return("", errors.New("token error"))
			},
			requestBody: "username=u&email=e@example.com&password=p&confirm-password=p",

			context:        true,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "",
		},
		{
			name: "db create token error",
			setupMocks: func() {
				resetter()

				form := formatters.RegisterForm{
					Username: "u",
					Email:    "e@example.com",
					Password: "p",
				}
				user := uschema.New().WithUsername(form.Username).WithEmail(form.Email).WithPassword(form.Password)

				mockErrResp.On("ErrorResponse", mock.Anything, mock.Anything, http.StatusInternalServerError).Return()
				mockUR.On("Create", mock.Anything, user).Return(&uschema.User{ID: 1, Username: "u", Email: "e"}, nil)
				mockTokens.On("NewToken", mock.Anything, mock.Anything).Return("token123", nil)
				mockTokR.On("CreateToken", mock.Anything, 1, "token123").Return(errors.New("token db error"))
			},
			requestBody: "username=u&email=e@example.com&password=p&confirm-password=p",

			context:        true,
			wantCode:       http.StatusInternalServerError,
			wantBodySubstr: "",
		},
		{
			name: "success registration",
			setupMocks: func() {
				resetter()

				template, err := template.New("welcome").Parse("<div>welcome</div>")
				assert.NoError(t, err)

				form := formatters.RegisterForm{Username: "u", Email: "e@example.com", Password: "p", Confirm: "p"}
				user := uschema.New().WithUsername(form.Username).WithEmail(form.Email).WithPassword(form.Password)

				mockUR.On("Create", mock.Anything, user).Return(&uschema.User{ID: 1, Username: "u", Email: "e@example.com"}, nil)
				mockTokens.On("NewToken", mock.Anything, mock.Anything).Return("token123", nil)
				mockTokR.On("CreateToken", mock.Anything, 1, "token123").Return(nil)
				mockVerificationEmail.On("SendVerificationEmail", form, "token123").Return()
				mockTemplate.On("Must", mock.Anything, mock.Anything).Return(template)
				mockTemplate.On("ParseFiles", mock.Anything).Return(template, nil)
				mockParser.On(
					"NewAlert",
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return("Success!!!")
			},
			requestBody: "username=u&email=e@example.com&password=p&confirm-password=p",

			context:        true,
			wantCode:       http.StatusOK,
			wantBodySubstr: "",
			wantPushUrl:    "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if tt.context {
				req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, reqCont))
			}

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
			if tt.wantBodySubstr != "" {
				assert.Contains(t, rec.Body.String(), tt.wantBodySubstr)
			}
			if tt.wantPushUrl != "" || tt.wantReswap != "" {
				assert.Equal(t, tt.wantPushUrl, rec.Header().Get("HX-Push-Url"))
				if tt.wantReswap != "" {
					assert.Equal(t, tt.wantReswap, rec.Header().Get("HX-Reswap"))
				}
			}
		})
	}
}
