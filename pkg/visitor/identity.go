package visitor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var anonymousIDPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func DeriveKey(secret, anonymousID, clientIP string) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", errors.New("visitor hash secret is required")
	}

	identity := ""
	if anonymousIDPattern.MatchString(strings.TrimSpace(anonymousID)) {
		identity = "anonymous:" + strings.ToLower(strings.TrimSpace(anonymousID))
	} else if strings.TrimSpace(clientIP) != "" {
		identity = "ip:" + strings.TrimSpace(clientIP)
	}
	if identity == "" {
		return "", errors.New("visitor identity is required")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(identity))
	return hex.EncodeToString(mac.Sum(nil)), nil
}
