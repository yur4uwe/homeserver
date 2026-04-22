package crypto

import (
	"homeserver/constants"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Getenv = os.Getenv
)

type CustomClaims struct {
	Seed string `json:"seed"`
	jwt.RegisteredClaims
}

func NewToken(seed string, expires time.Time) (string, error) {
	secret := []byte(Getenv(constants.TokenSecret))

	claims := &CustomClaims{
		Seed: seed,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenStr string) (*CustomClaims, bool) {
	if tokenStr == "" {
		return nil, false
	}

	secret := []byte(Getenv(constants.TokenSecret))

	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, false
	}

	return claims, token.Valid
}
