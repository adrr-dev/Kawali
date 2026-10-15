# PROGRESS.md — state sprint aktif MEWALI

## Sprint 1: fondasi + scaffolding server

- [x] Dokumen: `DESIGN.md`, `PRD.md`, `ARCHITECTURE.md`, `AGENTS.md` (§1-10)
- [x] Putusan MVP: mock+gray-out payment/shipping, anotasi cacat di foto, admin minimum, dua mode order
- [ ] Scaffold Go module + server runnable (in progress)
- [x] Scaffold Go module + server runnable (`go build`, `go vet`, smoke `/health` OK)
- [x] Makefile + unit test awal (`make test` hijau: config, db)
- [x] Model GORM + migrasi + seed filler (smoke: 3 akun terisi)
- [x] Logika murni: kode konfirmasi 5-char + ledger poin (test hijau)
- [x] Size-match toleransi 10% + skor kondisi aturan (putusan user, test hijau)
- [x] Layout templ + token CSS + Beranda skeleton (render + CSS terverifikasi)
- [x] Auth filler: sesi HMAC + login/logout + gate (unit + smoke OK)
- [x] Service katalog: draft + 6 slot + defect + publish gating + skor (test hijau)
- [x] Service order: cart + checkout + konfirmasi dua arah + poin done (11 test hijau)
- [x] Service review/lapor/chat/geo/ledger-redeem (semua paket hijau)
- [x] Halaman cari + detail + MatchBar live (smoke 87% terverifikasi)
- [x] Tambah-katalog: draft + upload 6 slot + marker klik + publish (smoke penuh OK)
- [x] Transaksi: keranjang/beli/checkout/konfirmasi/sukses + chat (smoke dua arah OK)
- [x] Poin/redeem + profil/ulasan/lapor + preferensi/lokasi (smoke OK)
- [x] Peta statis + admin skeleton moderasi (smoke: 403/200 OK)
- [x] Signup/lupa netral/onboarding + feed Beranda + chat di detail (smoke OK)
- [x] Final: Tailwind build + README + ARCHITECTURE as-implemented + SQL snapshot
- [ ] Model GORM inti + migrasi dev + seed filler
- [ ] Layout templ + token Tailwind + halaman awal

## Next steps

1. `go mod init` + deps (gin, gorm, sqlite/postgres, templ runtime)
2. `cmd/server`, `internal/config`, `internal/db`
3. `go build ./...` hijau
