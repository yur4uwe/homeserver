package crypto_test

import (
	"homeserver/constants"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks/mocksos"
	"homeserver/web-server/pkg/crypto"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	const testSecret = "test-secret"

	mockOS := new(mocksos.MockOS)
	mockOS.On("Getenv", constants.TokenSecret).Return(testSecret, nil)

	defer helpers.MockFunction(&crypto.Getenv, mockOS.Getenv, func() {
		mockOS.AssertExpectations(t)
	})()

	seed := "test-seed"

	// Generate a valid token
	tokenStr, err := crypto.NewToken(seed, time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// Validate a correct token
	claims, valid := crypto.ValidateToken(tokenStr)
	assert.True(t, valid, "Valid token should be accepted")
	assert.NotNil(t, claims)
	assert.Equal(t, seed, claims.Seed)

	// Tampered token (change a character)
	tampered := tokenStr[:len(tokenStr)-1]
	_, valid = crypto.ValidateToken(tampered)
	assert.False(t, valid, "Tampered token should not be valid")

	// Expired token
	expiredClaims := &crypto.CustomClaims{
		Seed: seed,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredStr, _ := expiredToken.SignedString([]byte(testSecret))
	_, valid = crypto.ValidateToken(expiredStr)
	assert.False(t, valid, "Expired token should not be valid")

	// Completely invalid token string
	_, valid = crypto.ValidateToken("not.a.jwt.token")
	assert.False(t, valid, "Completely invalid token string should not be valid")

	// Empty token string
	_, valid = crypto.ValidateToken("")
	assert.False(t, valid, "Empty token string should not be valid")

	// Token with wrong secret
	mockOS.ExpectedCalls = nil
	mockOS.On("Getenv", constants.TokenSecret).Return("wrong-secret").Once()
	_, valid = crypto.ValidateToken(tokenStr)
	assert.False(t, valid, "Token signed with different secret should not be valid")
	os.Setenv(constants.TokenSecret, testSecret) // restore
}
