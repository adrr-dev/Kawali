package components

import (
	"testing"
)

func TestPriceFormats(t *testing.T) {
	if got := PriceID(65000); got != "Rp65.000" {
		t.Errorf("PriceID = %q", got)
	}
	if got := PriceID(1500000); got != "Rp1.500.000" {
		t.Errorf("PriceID = %q", got)
	}
	if got := ShortPrice(65000); got != "65K" {
		t.Errorf("ShortPrice = %q", got)
	}
	if got := ShortPrice(65500); got != "Rp65.500" {
		t.Errorf("ShortPrice non-bulat = %q", got)
	}
}
