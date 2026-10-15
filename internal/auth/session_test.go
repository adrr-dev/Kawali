package auth

import (
	"testing"
	"time"
)

func TestSessionRoundtrip(t *testing.T) {
	now := time.Now()
	c := Issue("secret", 42, 30*24*time.Hour, now)
	uid, ok := Verify("secret", c, now.Add(time.Hour))
	if !ok || uid != 42 {
		t.Errorf("roundtrip gagal: %d,%v", uid, ok)
	}
}

func TestSessionTamperAndExpiry(t *testing.T) {
	now := time.Now()
	c := Issue("secret", 7, time.Hour, now)

	if _, ok := Verify("secret", c+"x", now); ok {
		t.Error("cookie rusak harus ditolak")
	}
	if _, ok := Verify("salah", c, now); ok {
		t.Error("secret salah harus ditolak")
	}
	if _, ok := Verify("secret", c, now.Add(2*time.Hour)); ok {
		t.Error("kedaluwarsa harus ditolak")
	}
	if _, ok := Verify("secret", "bukan-cookie", now); ok {
		t.Error("format salah harus ditolak")
	}
}
