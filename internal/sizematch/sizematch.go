package sizematch

import (
	"math"
)

// Margin toleransi: selisih 10% dari ukuran pengguna dianggap masih bisa
// ditampilkan (skor 0 untuk dimensi itu). Putusan user, lihat DECISIONS ADR-007.
const Margin = 0.10

type Sizes struct {
	Dada     float64
	Pinggang float64
	Pinggul  float64
	Tinggi   float64
	Torso    float64 // dipakai sebagai "panjang" barang
}

// Dim adalah satu baris perbandingan untuk UI ("Dada barang / kamu — 54/52").
type Dim struct {
	Name string
	Item float64
	User float64
}

type Result struct {
	Score    int   // 0-100
	Complete bool  // false bila Dada/Torso tak lengkap (lihat DESIGN Bab detail)
	Dims     []Dim // perbandingan per dimensi yang datanya ada
}

// Score menghitung kecocokan. Minimal Dada + Torso kedua sisi > 0 agar
// Complete; dimensi lain ikut bila datanya ada. Tinggi pengguna tidak
// dinilai (pengaruh ke fit kecil), hanya Dada/Pinggang/Pinggul/Torso.
func Score(user, item Sizes) Result {
	res := Result{}

	var sum float64
	var n int
	var dadaOK, torsoOK bool
	add := func(name string, u, it float64) {
		if u <= 0 || it <= 0 {
			return
		}
		res.Dims = append(res.Dims, Dim{Name: name, Item: it, User: u})
		s := 1 - math.Abs(it-u)/(Margin*u)
		sum += math.Max(0, math.Min(1, s))
		n++
		switch name {
		case "Dada":
			dadaOK = true
		case "Panjang":
			torsoOK = true
		}
	}
	add("Dada", user.Dada, item.Dada)
	add("Pinggang", user.Pinggang, item.Pinggang)
	add("Pinggul", user.Pinggul, item.Pinggul)
	add("Panjang", user.Torso, item.Torso)

	res.Complete = dadaOK && torsoOK
	if res.Complete && n > 0 {
		res.Score = int(math.Round(sum / float64(n) * 100))
	}
	return res
}
