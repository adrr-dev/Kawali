package components

import (
	"strconv"
	"strings"
)

// PriceID memformat rupiah penuh: 65000 -> "Rp65.000".
func PriceID(n int) string {
	s := strconv.Itoa(n)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	return "Rp" + strings.Join(parts, ".")
}

// ShortPrice meringkas harga kartu: 65000 -> "65K", selain itu harga penuh.
func ShortPrice(n int) string {
	if n >= 1000 && n%1000 == 0 {
		return strconv.Itoa(n/1000) + "K"
	}
	return PriceID(n)
}

// Utoa konversi uint untuk URL.
func Utoa(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// Itoa konversi int untuk label bernomor.
func Itoa(n int) string {
	if n < 0 {
		return "0"
	}
	return strconv.Itoa(n)
}

// Itoa64 konversi int64 untuk statistik.
func Itoa64(n int64) string {
	if n < 0 {
		return "0"
	}
	return strconv.FormatInt(n, 10)
}

// At menampilkan handle dengan @: At("rudi") -> "@rudi".
func At(handle string) string {
	return "@" + handle
}
