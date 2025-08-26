package crypto

import (
	"homeserver/constants"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	Getenv = os.Getenv
)

type CustomClaims struct {
	Seed string `json:"seed"`
	jwt.StandardClaims
}

func NewToken(seed string, expires time.Time) (string, error) {
	secret := []byte(Getenv(constants.TokenSecret))

	claims := &CustomClaims{
		Seed: seed,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
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
