package catalog

import (
	"testing"
)

func TestScoreTable(t *testing.T) {
	cases := []struct {
		photos, defects int
		want            int
		ok              bool
	}{
		{6, 0, 100, true},
		{6, 1, 90, true},
		{6, 2, 80, true},
		{6, 6, 40, true},  // lantai
		{6, 20, 40, true}, // lantai
		{7, 0, 100, true}, // foto lebih tetap ok
		{5, 0, 0, false},
		{0, 0, 0, false},
	}
	for _, c := range cases {
		got, ok := Score(c.photos, c.defects)
		if got != c.want || ok != c.ok {
			t.Errorf("Score(%d,%d) = %d,%v; mau %d,%v",
				c.photos, c.defects, got, ok, c.want, c.ok)
		}
	}
}
