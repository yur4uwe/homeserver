package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"homeserver/constants"
	"net/url"
	"strconv"
	"time"
)

func CheckHash(data, hash string) bool {
	return hmac.Equal([]byte(HashSeed(data)), []byte(hash))
}

func HashSeed(data string) string {
	secret := []byte(Getenv(constants.SourceSecret))
	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func SignedURL(rawurl string, expires time.Time) (string, error) {
	if expires.IsZero() {
		return "", fmt.Errorf("invalid input: expires must not be zero")
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	exp := expires.Unix()
	data := u.Path + "|" + strconv.FormatInt(exp, 10)

	signature := HashSeed(data)

	q := u.Query()
	q.Set("sid", signature)
	q.Set("exp", strconv.FormatInt(exp, 10))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func ValidateSignedURL(rawurl string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}

	q := u.Query()
	sid := q.Get("sid")
	expStr := q.Get("exp")
	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return false
	}

	data := u.Path + "|" + expStr
	return CheckHash(data, sid)
}
