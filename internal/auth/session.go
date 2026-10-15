package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Sesi cookie bertanda-tangan HMAC, tanpa dep tambahan:
// value = base64(userID) . expiryUnix . base64(sig)
const CookieName = "mewali_session"

// Issue membuat nilai cookie untuk userID dengan TTL dari now.
func Issue(secret string, userID uint, ttl time.Duration, now time.Time) string {
	uid := base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(userID), 10)))
	exp := strconv.FormatInt(now.Add(ttl).Unix(), 10)
	sig := sign(secret, uid, exp)
	return strings.Join([]string{uid, exp, sig}, ".")
}

// Verify memeriksa cookie. Mengembalikan userID bila valid dan belum kedaluwarsa.
func Verify(secret, cookie string, now time.Time) (uint, bool) {
	parts := strings.Split(cookie, ".")
	if len(parts) != 3 {
		return 0, false
	}
	uid, exp, sig := parts[0], parts[1], parts[2]
	if !hmac.Equal([]byte(sig), []byte(sign(secret, uid, exp))) {
		return 0, false
	}
	expUnix, err := strconv.ParseInt(exp, 10, 64)
	if err != nil || now.Unix() > expUnix {
		return 0, false
	}
	rawUID, err := base64.RawURLEncoding.DecodeString(uid)
	if err != nil {
		return 0, false
	}
	n, err := strconv.ParseUint(string(rawUID), 10, 64)
	if err != nil || n == 0 {
		return 0, false
	}
	return uint(n), true
}

func sign(secret, uid, exp string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "%s.%s", uid, exp)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// DefaultTTL sesi filler: 30 hari.
const DefaultTTL = 30 * 24 * time.Hour

// NewSession membuat nilai cookie + maxAge detik untuk login/signup.
func NewSession(secret string, userID uint, now time.Time) (value string, maxAge int) {
	return Issue(secret, userID, DefaultTTL, now), int(DefaultTTL.Seconds())
}
