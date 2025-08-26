package crypto_test

import (
	"homeserver/constants"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks/mocksos"
	"homeserver/web-server/pkg/crypto"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateSignedURL(t *testing.T) {
	mockOS := new(mocksos.MockOS)
	mockOS.On("Getenv", constants.SourceSecret).Return("secret", nil)

	defer helpers.MockFunction(&crypto.Getenv, mockOS.Getenv, func() {
		mockOS.AssertExpectations(t)
	})()

	validFileURL, _ := crypto.SignedURL("http://example.com/file/789", time.Now().Add(1*time.Hour))
	expiredFileURL, _ := crypto.SignedURL("http://example.com/file/999", time.Now().Add(-1*time.Hour))

	t.Run("valid signed url", func(t *testing.T) {
		assert.True(t, crypto.ValidateSignedURL(validFileURL))
	})

	t.Run("expired signed url", func(t *testing.T) {
		assert.False(t, crypto.ValidateSignedURL(expiredFileURL))
	})

	t.Run("tampered id", func(t *testing.T) {
		u := validFileURL
		u = u[:len(u)-1] + "0" // change last char of id
		assert.False(t, crypto.ValidateSignedURL(u))
	})

	t.Run("tampered entry type", func(t *testing.T) {
		u := validFileURL
		u = strings.Replace(u, "/file/", "/folder/", 1)
		assert.False(t, crypto.ValidateSignedURL(u))
	})

	t.Run("tampered signature", func(t *testing.T) {
		u := validFileURL
		// Replace sig param with an invalid value
		if idx := len(u) - 1; idx > 0 {
			u = u[:idx] + "0"
		}
		assert.False(t, crypto.ValidateSignedURL(u))
	})

	t.Run("missing signature", func(t *testing.T) {
		u := validFileURL
		// Remove sig param
		u = removeQueryParam(u, "sid")
		assert.False(t, crypto.ValidateSignedURL(u))
	})

	t.Run("missing expiry", func(t *testing.T) {
		u := validFileURL
		u = removeQueryParam(u, "exp")
		assert.False(t, crypto.ValidateSignedURL(u))
	})

	t.Run("missing id", func(t *testing.T) {
		u := validFileURL
		u = removePathPart(u, "789")
		assert.False(t, crypto.ValidateSignedURL(u))
	})

	t.Run("missing entry type", func(t *testing.T) {
		u := validFileURL
		u = removePathPart(u, "file")
		assert.False(t, crypto.ValidateSignedURL(u))
	})

	t.Run("invalid url", func(t *testing.T) {
		assert.False(t, crypto.ValidateSignedURL("::::not a url::::"))
	})
}

// Helper to remove a query param from a URL string
func removeQueryParam(rawurl, param string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl
	}
	q := u.Query()
	q.Del(param)
	u.RawQuery = q.Encode()
	return u.String()
}

func removePathPart(rawurl, part string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl
	}

	slice := strings.Split(u.Path, "/")
	for i := range slice {
		if slice[i] == part {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	u.Path = strings.Join(slice, "/")
	return u.String()
}
