package middleware_test

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"homeserver/web-server/pkg/crypto"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/templates"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthMiddleware(t *testing.T) {
	goodClaims := &crypto.CustomClaims{
		Seed: "test_seed",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(999999999),
		},
	}

	_, _, _, mockTokR, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, nil, nil, mockTokR, nil)
	mockTokR.On("QueryToken", mock.Anything, "").Return((*uschema.User)(nil), fmt.Errorf("token not found"))
	mockTokR.On("QueryToken", mock.Anything, "validToken").Return(uschema.NewUserBuilder(), nil)

	mockLogHistory := new(mocks.MockLogHistory)
	mockLogHistory.On("Log", mock.Anything).Return()
	mockLogHistory.On("ErrLog", mock.Anything, mock.Anything).Return()

	mockParseAsString := new(MockParseAsString)
	mockParseAsString.On("ParseTemplateToString", constants.Rootd+"static/templates/header.html", mock.Anything).Return("<header>Welcome</header>", nil)

	mockErrorResponse := new(mocks.ErrorResponseMock)
	mockErrorResponse.On("ErrorResponse", mock.Anything, mock.Anything, mock.Anything).Return()

	mockProtectedPath := new(MockIsProtectedPath)
	mockProtectedPath.On("IsProtectedPath", "/protected").Return(true)
	mockProtectedPath.On("IsProtectedPath", "/").Return(false)

	mockTokens := new(mocks.MockTokens)
	mockTokens.On("ValidateToken", "validToken").Return(goodClaims, true)
	mockTokens.On("ValidateToken", "").Return((*crypto.CustomClaims)(nil), false)

	cleanupParseAsString := helpers.MockFunction(&middleware.ParseTemplateToString, mockParseAsString.ParseTemplateToString, func() {})
	defer cleanupParseAsString()

	cleanupErrorResponse := helpers.MockFunction(&middleware.ErrorResponse, mockErrorResponse.ErrorResponse, func() {})
	defer cleanupErrorResponse()

	cleanupIsProtectedPath := helpers.MockFunction(&middleware.IsProtectedPath, mockProtectedPath.IsProtectedPath, func() {})
	defer cleanupIsProtectedPath()

	cleanupValidateToken := helpers.MockFunction(&middleware.ValidateToken, mockTokens.ValidateToken, func() {})
	defer cleanupValidateToken()

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mockRequestContext := &middleware.AdditionalRequestContext{
		LogHistory:  mockLogHistory,
		PageContent: &templates.LayoutPageData{},
	}

	t.Run("Handles missing request context", func(t *testing.T) {
		handler := middleware.AuthMiddleware(mockDB)(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Response status code should be 500 Internal Server Error")
		mockErrorResponse.AssertCalled(t, "ErrorResponse", w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
	})

	t.Run("Handles invalid or expired token", func(t *testing.T) {
		handler := middleware.AuthMiddleware(mockDB)(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "authTokenHS",
			Value: "",
		})
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockRequestContext))

		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		mockTokR.AssertCalled(t, "QueryToken", mock.Anything, "")
	})

	t.Run("Handles valid token", func(t *testing.T) {
		handler := middleware.AuthMiddleware(mockDB)(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "authTokenHS",
			Value: "validToken",
		})
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockRequestContext))

		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		mockTokR.AssertCalled(t, "QueryToken", mock.Anything, "validToken")
	})

	t.Run("Redirects to login for protected path without valid token", func(t *testing.T) {
		handler := middleware.AuthMiddleware(mockDB)(mockHandler)

		req := httptest.NewRequest("GET", "/protected", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockRequestContext))

		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusSeeOther, w.Code, "Response status code should be 303 See Other")
		assert.Equal(t, "/login", w.Header().Get("Location"), "Should redirect to /login")
	})

	t.Run("Handles template parsing error", func(t *testing.T) {
		mockErrorParseAsString := new(MockParseAsString)
		mockErrorParseAsString.On("ParseTemplateToString", constants.Rootd+"static/templates/header.html", mock.Anything).Return("", fmt.Errorf("template error"))
		cleanupParseAsString := helpers.MockFunction(&middleware.ParseTemplateToString, mockErrorParseAsString.ParseTemplateToString, func() {
			helpers.MockFunction(&middleware.ParseTemplateToString, mockParseAsString.ParseTemplateToString, func() {})
		})
		defer cleanupParseAsString()

		handler := middleware.AuthMiddleware(mockDB)(mockHandler)

		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockRequestContext))

		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Response status code should be 500 Internal Server Error")
		mockErrorResponse.AssertCalled(t, "ErrorResponse", w, "Internal Server Error: template error", http.StatusInternalServerError)
	})

	t.Run("Handles static file requests", func(t *testing.T) {
		mockProtectedPath.On("IsProtectedPath", "/static/").Return(false)

		handler := middleware.AuthMiddleware(mockDB)(mockHandler)

		req := httptest.NewRequest("GET", "/static/file.css", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ReqContKey, mockRequestContext))

		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
		mockProtectedPath.AssertNotCalled(t, "IsProtectedPath", "/static/file.css")
	})

	// mockTokR.AssertExpectations(t)
	// mockLogHistory.AssertExpectations(t)
	// mockParseAsString.AssertExpectations(t)
	// mockErrorResponse.AssertExpectations(t)
	// mockTokens.AssertExpectations(t)
}
