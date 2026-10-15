# MEWALI — NEW LIFE, SAME STYLE

Platform jual-beli pakaian thrift/preloved lokal (satu akun pembeli + penjual).
Stack: Go + Gin + a-h/templ + Tailwind v4 + HTMX + Alpine.js. DB GORM
(`gorm.Model`): SQLite lokal, Postgres Supabase saat deploy (ganti env saja).

Dokumen: `PRD.md` (cakupan) > PDF desain (visual) > `DESIGN.md` (pola) >
`ARCHITECTURE.md` (teknis) > `AGENTS.md` (aturan agent) >
`PROGRESS.md` / `DECISIONS.md` (state + ADR).

## Prasyarat

- Go 1.22+ (repo dibangun dengan 1.27), gcc (driver SQLite), `templ` CLI.
- Node 18+ untuk `make css`: `npm i -g @tailwindcss/cli tailwindcss`.
- Tidak perlu: Postgres, Maps API, gateway, OAuth (semua mock/filler di MVP).

## Jalan lokal (3 perintah)

```bash
cp .env.example .env   # opsional; default sudah jalan
make build             # templ generate + go build -> bin/mewali-server
./bin/mewali-server    # atau: make run
```

Buka `http://localhost:8080`. Akun filler (sandi `password`):
`pembeli@test.id`, `penjual@test.id`, `admin@test.id` (`/admin`).

Perintah lain: `make test` (wajib hijau), `make vet`, `make css`
(build `web/static/css/app.css` dari `input.css`), `make tidy`, `make clean`.

## Environment variables

| Nama | Default | Arti |
|---|---|---|
| `PORT` | `8080` | port listen |
| `DB_DRIVER` | `sqlite` | `sqlite` dev / `postgres` prod |
| `DATABASE_URL` | `file:./dev.db` | path sqlite atau URL postgres |
| `SESSION_SECRET` | `dev-secret…` | **wajib ganti di prod** (HMAC sesi) |
| `UPLOAD_DIR` | `./web/static/uploads` | penyimpanan foto lokal |
| `AUTO_MIGRATE` | `true` | `true` = AutoMigrate + seed saat start (dev). Prod: `false` + `migrations/*.sql` |
| `FEATURE_SHIPPING/GATEWAY/OAUTH` | `off` | cadangan flag (UI kini gray-out statis) |

Deploy Supabase: `DB_DRIVER=postgres`,
`DATABASE_URL=postgres://postgres:PASS@db.XXX.supabase.co:5432/postgres?sslmode=require`.
Tanpa ubah kode (satu `internal/db.Open()` switch dialek; larang SQL
dialek-spesifik). `AUTO_MIGRATE=false` di prod.

## Alur utama (coba manual)

1. Login penjual → `/tambah` (draft + ukuran) → `/katalog/:id`:
   unggah 6 slot, klik foto untuk penanda cacat, **Cek Rating & Terbitkan**.
2. Login pembeli → `/` / `/cari` → `/produk/:id` (skor, MatchBar bila
   preferensi terisi) → keranjang → **Beli** (pilih Cash/Transfer +
   Ketemuan-kode atau Chat-P2P) → bayar.
3. Meetup: `/konfirmasi/:id` dua sisi tukar kode 5-char (60 mnt) →
   `/sukses/:id` +10 poin. P2P: `/transaksi/:id` → tandai selesai manual.
4. `/poin` tukar reward (validasi saldo), `/profil/:handle` ulas + lapor,
   `/peta` jarak, `/chat` polling 5 dtk, `/admin` moderasi (admin saja).

## Batasan mock MVP (jujur, bukan bug)

- Bayar = Cash / Transfer-manual tercatat; QRIS/e-wallet/shipping = tombol
  disabled "Fitur belum tersedia". Tanpa webhook/escrow.
- Peta statis (tanpa vendor/GPS); jarak Haversine dari koordinat manual.
- Skor kondisi = aturan kelengkapan (`100-10×cacat`, lantai 40), bukan AI.
- Size-match = toleransi 10% (butuh Dada+Torso), bukan standar industri.
- Chat polling, bukan WebSocket. Reset sandi/OAuth Google-Apple belum ada.
- Suspend/blokir admin, level selain "Pemula Hijau", retur, sengketa
  kedaluwarsa: TBD produk (lihat `DESIGN.md` Bab 15).

## Testing

`make test` — 17 paket: config, db, models (17 tabel + unique),
auth (sesi HMAC, login, gate), catalog (slot/publish/skor, search),
orders (kode, checkout, konfirmasi dua arah, p2p, cart), points
(earn/redeem/guard), sizematch (+regresi 4-dim), geo, chat, reviews,
reports, storage (sniffing, 5MB), users (pref, register), seed
(idempoten), admin (resolve/takedown/audit), components (format harga).
Hermetis: SQLite `t.TempDir()`, tanpa layanan luar.
