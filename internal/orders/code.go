package orders

import (
	"crypto/rand"
	"strings"
	"time"
)

// Charset tanpa karakter ambigu O/0 dan I/1 (lihat DESIGN Bab konfirmasi).
const charset = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"

const (
	CodeLength = 5
	CodeTTL    = 60 * time.Minute
)

// Generate membuat kode acak sepanjang CodeLength dari charset.
func Generate() (string, error) {
	var sb strings.Builder
	sb.Grow(CodeLength)
	buf := make([]byte, CodeLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for _, b := range buf {
		sb.WriteByte(charset[int(b)%len(charset)])
	}
	return sb.String(), nil
}

// ValidFormat memeriksa panjang dan keanggotaan charset (termasuk menolak
// karakter ambigu O/0/I/1).
func ValidFormat(code string) bool {
	if len(code) != CodeLength {
		return false
	}
	for _, c := range code {
		if !strings.ContainsRune(charset, c) {
			return false
		}
	}
	return true
}

// ExpiredAt menghitung batas waktu dari momen pembuatan.
func ExpiredAt(from time.Time) time.Time {
	return from.Add(CodeTTL)
}

// IsExpired true bila now melewati expiresAt.
func IsExpired(expiresAt, now time.Time) bool {
	return now.After(expiresAt)
}
