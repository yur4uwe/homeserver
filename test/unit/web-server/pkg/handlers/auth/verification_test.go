package auth_test

import (
	"bytes"
	"context"
	"errors"
	"homeserver/constants"
	"homeserver/db/repositories"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"homeserver/web-server/pkg/handlers/auth"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/templates"
	"homeserver/web-server/pkg/utils/formatters"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type VerFailureMock struct {
	mock.Mock
}

func (m *VerFailureMock) verFailed(w http.ResponseWriter, verification *templates.VerificationPageData, code int) {
	if verification.Error == "" {
		panic("verFailed called with empty error")
	} else if verification.Success {
		panic("verFailed called with successful verification")
	}

	m.Called(w, verification, code)
	w.WriteHeader(code)
	w.Write([]byte(verification.Error))
}

func TestVerificationHandler(t *testing.T) {
	_, _, _, mockTokR, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, nil, nil, mockTokR, nil)
	mockParser := new(mocks.MockParser)
	mockErrResp := new(mocks.ErrorResponseMock)
	mockVerFail := new(VerFailureMock)
	mockTemplate := new(mocks.MockTemplate)
	mockReqCont := middleware.DefaultContext

	cleanupNewAlert := helpers.MockFunction(&auth.NewAlert, mockParser.NewAlert, func() {})
	cleanupLoadIcon := helpers.MockFunction(&auth.LoadIcon, mockParser.LoadIcon, func() {})
	cleanupErrResp := helpers.MockFunction(&auth.ErrorResponse, mockErrResp.ErrorResponse, func() {})
	cleanupVerFailed := helpers.MockFunction(&auth.VerFailed, mockVerFail.verFailed, func() {})
	cleanupMust := helpers.MockFunction(&auth.Must, mockTemplate.Must, func() {})
	cleanupParseFiles := helpers.MockFunction(&auth.ParseFiles, mockTemplate.ParseFiles, func() {})

	defer func() {
		cleanupNewAlert()
		cleanupLoadIcon()
		cleanupErrResp()
		cleanupVerFailed()
		cleanupMust()
		cleanupParseFiles()
	}()

	resetter := helpers.NewMocksResetter(&mockTokR.Mock, &mockParser.Mock)

	handler := auth.VerificationHandler(mockDB)

	newTemplate := func(inner_text string) *template.Template {
		return template.Must(template.New("verification").Parse(inner_text))
	}

	t.Run("missing request context", func(t *testing.T) {
		resetter()

		mockErrResp.On(
			"ErrorResponse",
			mock.Anything, mock.Anything, mock.Anything,
		).Return("Missing request context")

		req := httptest.NewRequest(http.MethodGet, "/verification?token=abc123", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Missing request context")
	})

	t.Run("missing token", func(t *testing.T) {
		resetter()

		mockVerFail.On("verFailed", mock.Anything, mock.Anything, mock.Anything).Return()
		mockParser.On("LoadIcon", mock.Anything).Return("icon")

		req := httptest.NewRequest(http.MethodGet, "/verification", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockReqCont))

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Missing verification token")
		mockVerFail.AssertCalled(t, "verFailed", rec, mock.Anything, http.StatusBadRequest)
	})

	t.Run("invalid token", func(t *testing.T) {
		resetter()

		mockTokR.On(
			"VerifyToken",
			mock.Anything, "badtoken",
		).Return(errors.New("invalid token"))

		mockVerFail.On(
			"verFailed",
			mock.Anything, mock.Anything, http.StatusInternalServerError,
		).Return()

		mockParser.On("LoadIcon", mock.Anything).Return("icon")

		req := httptest.NewRequest(http.MethodGet, "/verification?token=badtoken", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockReqCont))

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to check verification token")
	})

	t.Run("success", func(t *testing.T) {
		resetter()

		mockTokR.On("VerifyToken", mock.Anything, "goodtoken").Return(nil)
		mockParser.On("LoadIcon", mock.Anything).Return("icon")
		mockParser.On(
			"NewAlert",
			"Verification success", "check",
			"Your email has been verified successfully!", constants.Success,
		).Return("ALERT: Verification success check Your email has been verified successfully! 1")

		mockTemplate.On(
			"Must",
			mock.Anything, mock.Anything,
		).Return(newTemplate("Your email has been verified successfully!"))

		mockTemplate.On(
			"ParseFiles",
			mock.Anything,
		).Return(newTemplate("<html><body>Verification</body></html>"), nil)

		req := httptest.NewRequest(http.MethodGet, "/verification?token=goodtoken", nil)
		rec := httptest.NewRecorder()

		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockReqCont))

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Your email has been verified")
	})
}

func TestSendVerificationFailurePage(t *testing.T) {
	content := "<html><body>{{.Error}}</body></html>"

	mockedTemplate := new(mocks.MockTemplate)

	defer helpers.MockFunction(&auth.ParseFiles, mockedTemplate.ParseFiles, func() {})()
	defer helpers.MockFunction(&auth.Must, mockedTemplate.Must, func() {})()

	mockedTemplate.On("ParseFiles", mock.Anything).
		Return(&template.Template{}, nil)

	mockedTemplate.On("Must", mock.Anything, mock.Anything).
		Return(template.Must(template.New("test").Parse(content)))

	rec := httptest.NewRecorder()
	verification := &templates.VerificationPageData{
		Success: false,
		Error:   "Test error message",
		Causes:  []template.HTML{"Cause 1", "Cause 2"},
		Icon:    "",
	}

	auth.SendVerificationFailurePage(rec, verification, http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test error message")
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
}

func TestSendVerificationEmail(t *testing.T) {
	// Set environment variables for email config
	os.Setenv(constants.EmailPort, "2525")
	os.Setenv(constants.EmailHost, "smtp.example.com")
	os.Setenv(constants.EmailUser, "user")
	os.Setenv(constants.EmailPassword, "pass")
	defer func() {
		os.Unsetenv(constants.EmailPort)
		os.Unsetenv(constants.EmailHost)
		os.Unsetenv(constants.EmailUser)
		os.Unsetenv(constants.EmailPassword)
	}()

	form := formatters.RegisterForm{
		Username: "testuser",
		Email:    "test@example.com",
	}

	// This test just ensures no panic and prints to stderr if email sending fails
	verificationToken := "abc123"
	// Redirect stderr to buffer to suppress output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	auth.SendVerificationEmail(form, verificationToken)

	w.Close()
	os.Stderr = oldStderr
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	output := buf.String()
	assert.Contains(t, output, "Failed to send verification email")
}
