package sizematch

import (
	"testing"
)

func TestExactMatch100(t *testing.T) {
	u := Sizes{Dada: 52, Pinggang: 80, Pinggul: 90, Tinggi: 170, Torso: 70}
	r := Score(u, u)
	if !r.Complete || r.Score != 100 {
		t.Errorf("identik harus 100+complete, dapat %+v", r)
	}
}

func TestTenPercentOffScoresZero(t *testing.T) {
	u := Sizes{Dada: 50, Torso: 70}
	it := Sizes{Dada: 55, Torso: 77} // +10% di semua dimensi
	r := Score(u, it)
	if !r.Complete {
		t.Fatal("harus complete")
	}
	if r.Score != 0 {
		t.Errorf("skor = %d, mau 0", r.Score)
	}
}

func TestFivePercentOffScoresFifty(t *testing.T) {
	u := Sizes{Dada: 50, Torso: 70}
	it := Sizes{Dada: 52.5, Torso: 73.5} // +5% = tengah toleransi
	r := Score(u, it)
	if r.Score != 50 {
		t.Errorf("skor = %d, mau 50", r.Score)
	}
}

func TestIncompleteWithoutDadaOrTorso(t *testing.T) {
	u := Sizes{Dada: 50, Torso: 70}
	r := Score(u, Sizes{Dada: 50}) // torso barang kosong
	if r.Complete {
		t.Error("tanpa torso harus incomplete")
	}
	if r.Score != 0 {
		t.Errorf("incomplete harus skor 0, dapat %d", r.Score)
	}
	r2 := Score(Sizes{}, Sizes{Dada: 50, Torso: 70})
	if r2.Complete {
		t.Error("tanpa preferensi harus incomplete")
	}
}

func TestDimsForUI(t *testing.T) {
	u := Sizes{Dada: 52, Torso: 70}
	r := Score(u, Sizes{Dada: 54, Torso: 71})
	if len(r.Dims) != 2 {
		t.Fatalf("dims = %d, mau 2", len(r.Dims))
	}
	if r.Dims[0].Name != "Dada" || r.Dims[0].Item != 54 || r.Dims[0].User != 52 {
		t.Errorf("dim dada salah: %+v", r.Dims[0])
	}
}

func TestFourDimsAverage(t *testing.T) {
	u := Sizes{Dada: 52, Pinggang: 80, Pinggul: 90, Torso: 70}
	it := Sizes{Dada: 54, Pinggang: 80, Pinggul: 90, Torso: 71}
	r := Score(u, it)
	if !r.Complete {
		t.Fatal("harus complete")
	}
	if len(r.Dims) != 4 {
		t.Fatalf("dims = %d, mau 4", len(r.Dims))
	}
	if r.Score != 87 {
		t.Errorf("skor = %d, mau 87", r.Score)
	}
}
