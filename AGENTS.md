# AGENTS.md — aturan kerja agent di repo MEWALI

## 1. Urutan baca wajib

1. `PRD.md` — cakupan fitur (apa yang boleh dibangun).
2. PDF desain terbaru (visual, bila ada) + `DESIGN.md` — pola UI, token, state.
3. `ARCHITECTURE.md` — stack, struktur, alur bisnis.

Jika konflik: PRD > PDF > `DESIGN.md` > preferensi pribadi. Jangan mengarang layar/aturan bisnis yang ditandai TBD di `DESIGN.md:616`.

## 2. Stack (jangan diganti tanpa alasan)

Go + Gin + a-h/templ + Tailwind v4 + HTMX + Alpine/JS kecil. DB GORM + `gorm.Model`. SQLite dev, Postgres Supabase prod via env saja.

## 3. Perintah wajib

```bash
templ generate          # wajib sebelum go build/run
go run ./cmd/server
go test ./...           # sqlite file/memory
go vet ./...            # bila menyentuh backend
```

Tailwind v4 build terpisah (tanpa CDN di prod). Jangan commit `web/static/uploads/*`.

## 4. Gaya kode

- Lapisan: `handler -> service -> repo -> model`. Handler dilarang panggil GORM langsung.
- Semua model embed `gorm.Model`. Uang `int` rupiah, koordinat `double`. Larangan: SQL dialek-spesifik (`strftime`, FTS5, `AUTOINCREMENT` mentah, fungsi geo SQL).
- `sizematch`, charset kode, hitung poin: pure function + unit test.
- Env: `PORT, DB_DRIVER, DATABASE_URL, SESSION_SECRET, UPLOAD_DIR, FEATURE_*`. Tidak ada branch driver di repo/service.

## 5. Gaya UI

- Bahasa Indonesia, sapaan "kamu". Judul layar HURUF KAPITAL. Format `Rp65.000`, `65K` di kartu, desimal koma (`4,8`, `0,8 km`), ukuran `cm`.
- Token via `@theme` saja; larang hex di `.templ`. Tap >=44px, ikon-saja wajib `aria-label`.
- Setiap fetch: `Loading/Empty/Error/Success/Disabled` (`DESIGN.md:207`). Fitur belum tersedia = disabled + tooltip "Fitur belum tersedia", di belakang `feature` flag (`FEATURE_SHIPPING=off`, `FEATURE_GATEWAY=off`, `FEATURE_OAUTH=off`).
- Placeholder (Rusdi, `K9BS0`, `89%`, `4,8`) hanya di `seed/` fixture berlabel demo. Dilarang jadi default produksi.

## 6. Larangan klaim palsu

Dilarang render seolah aktif: payment gateway, shipping/tracking, skor AI, badge "terverifikasi/aman", rating `0,0`, poin sebelum konfirmasi backend, alamat rumah presisi ke publik.

## 7. Sebelum selesai

- [ ] `templ generate` + build/test hijau.
- [ ] Checklist `DESIGN.md:633` yang relevan terpenuhi.
- [ ] TBD baru ditandai eksplisit, bukan ditebak.

## 8. Git & secrets

- Agent boleh `git add`/`git commit` lokal (pesan konvensional), tapi **dilarang `git push`** — push dilakukan pemilik repo sendiri.
- Dilarang commit: `.env`, `*.db` / `*.sqlite*`, `web/static/uploads/*`. Contoh env hanya via `.env.example`.

## 9. Protokol building (berlaku saat kode asli mulai dibangun)

- Setiap selesai satu task: update `PROGRESS.md` (state sprint aktif, sub-task selesai, langkah berikut). Tujuannya agar state selamat melewati kompresi token.
- Setiap keputusan arsitektur saat eksekusi (pilih lib, ubah skema, trade-off): catat ADR ringkas di `DECISIONS.md` (konteks, opsi, keputusan, alasan).
- Kedua file dibuat saat fase building dimulai, bukan sekarang.

## 10. Stop-and-ask (TBD — dilarang menebak)

Wajib berhenti dan tanya user bila menyentuh: definisi/rentang skor kondisi + metode penilaian; rumus % size-match + ukuran minimum; ambang level/expired/retur poin; sengketa kode (salah/expired/tidak merespons); eligibility/edit/hapus ulasan; gateway/shipping/tracking real; provider peta + presisi lokasi publik; permission admin, suspend/banding; draft/edit/nonaktifkan listing.
