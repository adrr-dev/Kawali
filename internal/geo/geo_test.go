package geo

import (
	"math"
	"testing"
)

func TestHaversineEquator(t *testing.T) {
	// 1 derajat bujur di ekuator ~= 111,19 km
	got := HaversineKm(0, 0, 0, 1)
	if math.Abs(got-111.19) > 0.5 {
		t.Errorf("jarak = %v, mau ~111,19", got)
	}
	if d := HaversineKm(-6.2, 106.8, -6.2, 106.8); d != 0 {
		t.Errorf("sama harus 0, dapat %v", d)
	}
}

func TestFormatID(t *testing.T) {
	if s := FormatID(0.8); s != "0,8 km" {
		t.Errorf("format = %q", s)
	}
}
