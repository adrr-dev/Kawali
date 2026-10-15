package catalog

// Aturan skor kondisi MVP (putusan user, ADR-007): berbasis kelengkapan,
// tanpa klaim AI. 6 foto wajib agar bisa dinilai; tiap cacat -10 poin,
// lantai 40. Method selalu "aturan".
const (
	PhotosRequired = 6
	DefectPenalty  = 10
	ScoreFloor     = 40
	ScoreMax       = 100
	MethodRule     = "aturan"
)

// Score mengembalikan skor + true bila foto lengkap. Bila foto kurang,
// ok=false dan skor tidak boleh dipakai (UI tampilkan "belum dinilai").
func Score(photoCount, defectCount int) (score int, ok bool) {
	if photoCount < PhotosRequired {
		return 0, false
	}
	if defectCount < 0 {
		defectCount = 0
	}
	s := ScoreMax - DefectPenalty*defectCount
	if s < ScoreFloor {
		s = ScoreFloor
	}
	return s, true
}
