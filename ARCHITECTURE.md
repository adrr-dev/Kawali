# ARCHITECTURE.md — MEWALI (as implemented, MVP)

> Stack: Go + Gin + a-h/templ + Tailwind v4 + HTMX + Alpine.js. Modular monolith,
> server-rendered. `make build` hijau, `make test` 17 paket hijau (Oktober 2026).
> Prioritas sumber: **PRD** > PDF desain > `DESIGN.md` > file ini.
> Riwayat putusan: `DECISIONS.md` (ADR-001–017). State sprint: `PROGRESS.md`.

## 1. Cara jalan (detail di `README.md`)

```bash
make build && ./bin/mewali-server   # :8080, SQLite dev.db, seed filler
make test | make vet | make css
```

Dua permukaan: aplikasi pengguna mobile-first (satu akun buyer+seller) dan
`/admin` desktop untuk moderasi minimum. Flow selalu
`browser -> Gin -> service -> GORM -> DB`. Backend sumber kebenaran untuk
harga, poin, kode, status, moderasi.

## 2. Struktur file aktual

```
cmd/server/main.go            # wiring: config, db, migrate+seed, routes, listen
internal/
  admin/{service,handler}.go       # overview, resolve/reject laporan, takedown + audit
  auth/{session,lookup,service,handler}.go  # cookie HMAC, login/logout, gate
  catalog/{service,search,handler,manage}.go# draft, slot, defect, publish, cari, detail, kelola
  chat/{service,handler}.go        # thread reuse, posting, partial polling
  config/config.go                 # PORT, DB_DRIVER, DATABASE_URL, SESSION_SECRET, UPLOAD_DIR, AUTO_MIGRATE
  db/{db,migrate}.go               # Open() switch dialek (satu-satunya branch driver)
  geo/{geo,handler}.go             # Haversine, FormatID koma, /peta statis
  media/cover.go                   # CoverURL + fallback (rumah tunggal anti-siklus)
  models/models.go                 # 17 model + konstanta status + Categories
  orders/{code,service,handler}.go # kode 5-char, cart, checkout, konfirmasi, sukses
  points/{points,service,handler}.go# tabel earn, ledger append-only, redeem atomik
  reports/service.go               # FileReport (open)
  reviews/service.go               # eligibility + SellerRating
  seed/seed.go                     # filler + demo product/banner/donasi (idempoten)
  sizematch/sizematch.go           # toleransi 10%, wajib Dada+Torso
  storage/local.go                 # Storage iface + disk (sniffing, 5MB, nama acak)
  users/{pref,register,handler}.go # preferensi, lokasi, signup, profil, ulasan, lapor
web/templates/{layout,components,pages}/  # *.templ + *_templ.go (jangan edit manual)
web/static/{css/{input,app}.css,img/placeholder.svg,uploads/}
migrations/0001_init.sql           # snapshot skema dev (SQLite); review sebelum prod
Makefile (.env.example, .gitignore)
```

Aturan lapis (ADR-012/015/017): `handler -> service -> repo/fungsi -> model`.
Arah impor tunggal: handler→service, pages→components, users→auth
(lookup user terpusat di `auth`; `users/repo.go` dihapus karena siklus).
Helper bersama: `models` (konstanta), `components` (format), `media` (cover),
`seed` (paket sendiri agar boleh pakai service).

## 3. Endpoint aktual

| Method | Route | Akses | Keterangan |
|---|---|---|---|
| GET | `/`, `/cari`, `/produk/:id` | publik | feed, search+filter, detail+MatchBar |
| POST | `/cart` | login | tambah keranjang |
| GET/POST | `/login`, `/logout`, `/signup`, `/lupa-sandi`, `/onboarding`, `/health` | publik | auth filler + netral reset |
| GET/POST | `/tambah`, `/katalog`, `/katalog/:id`, `…/foto`, `…/cacat`, `…/publish`, `/barang-dijual` | seller pemilik | draft→6 slot→marker→publish |
| GET/POST | `/keranjang`, `/beli`, `/checkout` | login | cart, ringkasan+gate lokasi, buat order |
| GET/POST | `/konfirmasi/:id` | anggota order | kode lawan, timer Alpine, chat link |
| GET/POST | `/transaksi/:id`, `…/selesai` | anggota order | status P2P + selesai manual |
| GET | `/sukses/:id` | anggota order | done + poin ledger |
| GET/POST | `/poin`, `/redeem` | login | saldo, earn list, tukar (validasi) |
| GET | `/profil`, `/profil/:handle`, `/preferensi` | login | profil, ulasan, ukuran+lokasi |
| POST | `/preferensi`, `/ulasan`, `/lapor` | login | simpan, eligibility server |
| POST/GET | `/chat/mulai`, `/chat/:id`, `/partial/chat/:id` | anggota thread | mulai, isi+polling 5 dtk |
| GET | `/peta` | login | pin produk/donasi + jarak koma |
| GET/POST | `/admin`, `/admin/laporan`, `/admin/produk`, `/admin/transaksi` | admin | ringkasan, resolve/reject, takedown |

Semua halaman pakai `layout.Base` + `Dock`. Partial HTMX hanya
`/partial/chat/:id`. Upload: multipart → sniff `DetectContentType` →
`storage.Local` → URL.

## 4. Skema (semua embed `gorm.Model`)

`users(is_admin)`, `size_preferences(user unique)`, `products(status:
draft/in_review/active/sold/taken_down)`, `product_photos`,
`product_sizes`, `defects(x,y 0-1000)`, `condition_scores(product unique,
method)`, `cart_items(user+product unique)`, `orders(pay: cash/
transfer_manual; mode: p2p/meetup; status)`, `confirm_codes(order+role
unique, code 5, expires_at, used_at, attempts)`, `point_ledgers(append-only,
balance_after)`, `redemptions`, `reviews(order unique)`, `reports(status)`,
`chat_threads`, `chat_messages`, `donation_points`, `banners`,
`admin_audits`. Uang `int` rupiah, koordinat `double`. Indeks unik:
email, handle, cart, code, review, product_size, condition_score.
Snapshot: `migrations/0001_init.sql` (19 tabel).

## 5. DB & deploy

`DB_DRIVER=sqlite/postgres`, `DATABASE_URL` — hanya `internal/db.Open()`
yang branch. Larang SQL dialek-spesifik; geo di Go; search via
`LOWER(col) LIKE`. `AUTO_MIGRATE=true` dev (migrate+seed); prod `false` +
SQL berversion. Upload lokal dev; ganti `storage.Storage` ke bucket tanpa
ubah skema (kolom tetap URL string).

## 6. Aturan bisnis terimplementasi

- Publish: 6 slot + skor `100-10×cacat` (lantai 40, method `aturan`).
- Match: per dimensi `1-|i-u|/(0.10·u)`, rata-rata bulat; Complete wajib
  Dada+Torso; Tinggi tak dinilai (ADR-007).
- Kode: charset `ABCDEFGHJKMNPQRSTUVWXYZ23456789`, TTL 60 mnt, wajib kode
  lawan; done (2 sisi) → sold + ledger +10; guard satu pending/produk.
- Poin: earn +10/+15/+5/+20/+25; redeem 100/80/120 atomik; saldo dari ledger.
- Ulasan: done + buyer + 1/order; kosong → "Belum ada rating".
- Checkout diblokir bila lat+lng nol. Gateway/shipping/retur/tracking =
  disabled "Fitur belum tersedia". Poin tampil hanya pasca-done.
- Admin: tanpa suspend/blokir/hapus (TBD); takedown bukan hapus; tiap aksi
  tercatat `admin_audits` (belum ada UI daftar audit — follow-up).

## 7. Frontend

Token di `@theme` (`input.css`) + gaya komponen ikut dibuild ke `app.css`
(`make css`; butuh `@tailwindcss/cli` global). Komponen: Button/Field/
Chip/Card/Dock/PointsBadge/MatchBar/ChatBubble/Sheet-lite/Empty/Error +
admin table. Bahasa Indonesia, sapaan "kamu", judul kapital, `Rp65.000`,
`65K`, koma desimal, `cm`. Tap ≥44px, `aria-label` ikon-saja, focus ring,
safe-area. Placeholder hanya `seed` demo berlabel.
