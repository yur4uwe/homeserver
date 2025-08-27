package mocks

import (
	"homeserver/pkg/crypto"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockTokens struct {
	mock.Mock
}

func (m *MockTokens) NewToken(seed string, expires time.Time) (string, error) {
	args := m.Called(seed, expires)
	return args.String(0), args.Error(1)
}

func (m *MockTokens) ValidateToken(token string) (*crypto.CustomClaims, bool) {
	args := m.Called(token)
	if !args.Bool(1) {
		return nil, false
	}
	return args.Get(0).(*crypto.CustomClaims), true
}
