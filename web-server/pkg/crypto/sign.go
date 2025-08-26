package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"homeserver/constants"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func SignedURL(rawurl string, entryType, id string, expires time.Time) (string, error) {
	if entryType == "" || id == "" {
		return "", fmt.Errorf("invalid input: entryType and id must not be empty")
	}

	secret := []byte(Getenv(constants.SourceSecret))

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	exp := expires.Unix()
	data := entryType + "|" + id + "|" + strconv.FormatInt(exp, 10)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(data))
	signature := hex.EncodeToString(mac.Sum(nil))

	q := u.Query()
	q.Set("sid", signature)
	q.Set("exp", strconv.FormatInt(exp, 10))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func ValidateSignedURL(rawurl string) bool {
	secret := []byte(Getenv(constants.SourceSecret))

	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}

	parts := strings.Split(u.Path, "/")
	if len(parts) < 3 {
		return false
	}

	entryType := parts[1]
	id := parts[2]

	q := u.Query()
	sid := q.Get("sid")
	expStr := q.Get("exp")
	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return false
	}
	data := entryType + "|" + id + "|" + expStr
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(data))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(sid), []byte(expectedSig))
}
