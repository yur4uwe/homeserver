package auth

import (
	"bytes"
	"context"
	"errors"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/utils"
	"homeserver/test"
	"homeserver/test/mocks"
	"homeserver/web-server/pkg/templates"
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
	mockDB := new(mocks.MockDB)
	mockParser := new(mocks.MockParser)
	mockErrResp := new(mocks.ErrorResponseMock)
	mockVerFail := new(VerFailureMock)
	mockTemplate := new(mocks.MockTemplate)
	mockReqCont := middleware.DefaultContext

	cleanupNewAlert := test.MockFunction(&NewAlert, mockParser.NewAlert, func() {})
	cleanupLoadIcon := test.MockFunction(&LoadIcon, mockParser.LoadIcon, func() {})
	cleanupErrResp := test.MockFunction(&ErrorResponse, mockErrResp.ErrorResponse, func() {})
	cleanupVerFailed := test.MockFunction(&verFailed, mockVerFail.verFailed, func() {})
	cleanupMust := test.MockFunction(&Must, mockTemplate.Must, func() {})
	cleanupParseFiles := test.MockFunction(&ParseFiles, mockTemplate.ParseFiles, func() {})

	defer func() {
		cleanupNewAlert()
		cleanupLoadIcon()
		cleanupErrResp()
		cleanupVerFailed()
		cleanupMust()
		cleanupParseFiles()
	}()

	resetter := test.NewMocksResetter(&mockDB.Mock, &mockParser.Mock)

	handler := VerificationHandler(mockDB)

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

		mockDB.On(
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

		mockDB.On("VerifyToken", mock.Anything, "goodtoken").Return(nil)
		mockParser.On("LoadIcon", mock.Anything).Return("icon")
		mockParser.On(
			"NewAlert",
			"Verification success", "check",
			"Your email has been verified successfully!", utils.AlertSuccess,
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
	// Create a dummy template file
	tmpFile := "test_verification.html"
	content := "<html><body>{{.Error}}</body></html>"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpFile)

	// Patch ParseFiles and Must to use our dummy template
	origParseFiles := ParseFiles
	origMust := Must
	ParseFiles = func(filenames ...string) (*template.Template, error) {
		return template.ParseFiles(tmpFile)
	}
	Must = func(tmpl *template.Template, err error) *template.Template {
		return tmpl
	}
	defer func() {
		ParseFiles = origParseFiles
		Must = origMust
	}()

	rec := httptest.NewRecorder()
	verification := &templates.VerificationPageData{
		Success: false,
		Error:   "Test error message",
		Causes:  []template.HTML{"Cause 1", "Cause 2"},
		Icon:    "",
	}

	sendVerificationFailurePage(rec, verification, http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test error message")
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
}

func TestSendVerificationEmail(t *testing.T) {
	// Set environment variables for email config
	os.Setenv(globals.EnvEmailPort, "2525")
	os.Setenv(globals.EnvEmailHost, "smtp.example.com")
	os.Setenv(globals.EnvEmailUser, "user")
	os.Setenv(globals.EnvEmailPassword, "pass")
	defer func() {
		os.Unsetenv(globals.EnvEmailPort)
		os.Unsetenv(globals.EnvEmailHost)
		os.Unsetenv(globals.EnvEmailUser)
		os.Unsetenv(globals.EnvEmailPassword)
	}()

	form := globals.RegisterForm{
		Username: "testuser",
		Email:    "test@example.com",
	}

	// This test just ensures no panic and prints to stderr if email sending fails
	verificationToken := "abc123"
	// Redirect stderr to buffer to suppress output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	sendVerificationEmail(form, verificationToken)

	w.Close()
	os.Stderr = oldStderr
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	output := buf.String()
	assert.Contains(t, output, "Failed to send verification email")
}
