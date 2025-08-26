package crypto_test

import (
	"homeserver/constants"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks/mocksos"
	"homeserver/web-server/pkg/crypto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignedURL(t *testing.T) {
	mockOS := new(mocksos.MockOS)
	mockOS.On("Getenv", constants.SourceSecret).Return("secret", nil)

	defer helpers.MockFunction(&crypto.Getenv, mockOS.Getenv, func() {
		mockOS.AssertExpectations(t)
	})()

	tt := []struct {
		name      string
		rawurl    string
		entryType string
		id        string
		expires   time.Time
		isValid   bool
		shouldErr bool
	}{
		{
			name:      "valid url and time",
			rawurl:    "http://example.com/resource/123",
			entryType: "resource",
			id:        "123",
			expires:   time.Now().Add(1 * time.Hour),
			isValid:   true,
			shouldErr: false,
		},
		{
			name:      "expired url",
			rawurl:    "http://example.com/resource/456",
			entryType: "resource",
			id:        "456",
			expires:   time.Now().Add(-1 * time.Hour),
			isValid:   false,
			shouldErr: false,
		},
		{
			name:      "invalid url format",
			rawurl:    "::::not a url::::",
			entryType: "resource",
			id:        "789",
			expires:   time.Now().Add(1 * time.Hour),
			isValid:   false,
			shouldErr: true,
		},
		{
			name:      "empty id",
			rawurl:    "http://example.com/resource/",
			entryType: "resource",
			id:        "",
			expires:   time.Now().Add(1 * time.Hour),
			isValid:   false,
			shouldErr: true,
		},
		{
			name:      "empty entryType",
			rawurl:    "http://example.com//123",
			entryType: "",
			id:        "123",
			expires:   time.Now().Add(1 * time.Hour),
			isValid:   false,
			shouldErr: true,
		},
		{
			name:      "url with query params",
			rawurl:    "http://example.com/resource/123?foo=bar",
			entryType: "resource",
			id:        "123",
			expires:   time.Now().Add(1 * time.Hour),
			isValid:   true,
			shouldErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name+": "+tc.rawurl, func(t *testing.T) {
			signedURL, err := crypto.SignedURL(tc.rawurl, tc.entryType, tc.id, tc.expires)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.isValid, crypto.ValidateSignedURL(signedURL))
		})
	}
}
