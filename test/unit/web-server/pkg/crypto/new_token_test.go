package crypto_test

import (
	"os"
	"testing"
	"time"

	"homeserver/constants"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks/mocksos"
	"homeserver/web-server/pkg/crypto"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	// Set a test secret for token signing
	const testSecret = "test-secret"

	mockOS := new(mocksos.MockOS)
	mockOS.On("Getenv", constants.TokenSecret).Return(testSecret, nil)

	defer helpers.MockFunction(&crypto.Getenv, mockOS.Getenv, func() {
		mockOS.AssertExpectations(t)
	})()

	seed := "test-seed"

	tokenStr, err := crypto.NewToken(seed)
	assert.NoError(t, err, "NewToken should not return error")
	assert.NotEmpty(t, tokenStr, "NewToken should return a token string")

	// Parse the token to verify its claims
	token, err := jwt.ParseWithClaims(tokenStr, &crypto.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	assert.NoError(t, err, "Token should parse without error")
	assert.True(t, token.Valid, "Token should be valid")

	claims, ok := token.Claims.(*crypto.CustomClaims)
	assert.True(t, ok, "Claims should be of type CustomClaims")
	assert.Equal(t, seed, claims.Seed, "Seed claim should match input")

	// Check expiration is roughly 24h in the future
	now := time.Now().Unix()
	assert.InDelta(t, now+24*3600, claims.ExpiresAt, 5, "ExpiresAt should be about 24h in the future")
}

func TestNewToken_InvalidInput(t *testing.T) {
	const testSecret = "test-secret"
	os.Setenv(constants.TokenSecret, testSecret)
	defer os.Unsetenv(constants.TokenSecret)

	// Empty seed
	tokenStr, err := crypto.NewToken("")
	assert.NoError(t, err, "Should still create a token with empty seed")
	assert.NotEmpty(t, tokenStr)

	// Very long seed
	longSeed := make([]byte, 1024)
	for i := range longSeed {
		longSeed[i] = 'a'
	}
	tokenStr, err = crypto.NewToken(string(longSeed))
	assert.NoError(t, err, "Should create a token with a long seed")
	assert.NotEmpty(t, tokenStr)
}
