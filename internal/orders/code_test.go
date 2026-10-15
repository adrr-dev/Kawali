package orders

import (
	"testing"
	"time"
)

func TestGenerateLengthAndCharset(t *testing.T) {
	for i := 0; i < 200; i++ {
		c, err := Generate()
		if err != nil {
			t.Fatalf("Generate: %v", err)
		}
		if !ValidFormat(c) {
			t.Fatalf("kode %q lolos generate tapi gagal validasi", c)
		}
	}
}

func TestGenerateUniqueness(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 500; i++ {
		c, err := Generate()
		if err != nil {
			t.Fatalf("Generate: %v", err)
		}
		if seen[c] {
			t.Fatalf("kode duplikat %q dalam 500 generate", c)
		}
		seen[c] = true
	}
}

func TestValidFormatRejectsAmbiguous(t *testing.T) {
	for _, c := range []string{"K9BS0", "J3FO8", "ABO12", "AB01C", "ABCI1", "abcde", "ABCD", "ABCDEF", ""} {
		if ValidFormat(c) {
			t.Errorf("%q harus ditolak", c)
		}
	}
}

func TestValidFormatAccepts(t *testing.T) {
	for _, c := range []string{"K9BS2", "J3FZ8", "ABCDE", "23456"} {
		if !ValidFormat(c) {
			t.Errorf("%q harus diterima", c)
		}
	}
}

func TestExpiry(t *testing.T) {
	now := time.Now()
	exp := ExpiredAt(now)
	if exp.Sub(now) != CodeTTL {
		t.Errorf("TTL = %v, mau %v", exp.Sub(now), CodeTTL)
	}
	if IsExpired(exp, exp.Add(-time.Minute)) {
		t.Error("belum expired harus false")
	}
	if !IsExpired(exp, exp.Add(time.Second)) {
		t.Error("lewat batas harus true")
	}
}
