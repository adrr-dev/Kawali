# DECISIONS.md — Architectural Decision Records

## ADR-001: Module path `mewali`
- Konteks: repo bukan git remote, user push sendiri nanti.
- Keputusan: module `mewali` (tanpa host). Ganti saat remote final ada.
- Alasan: import lokal stabil, rename murah sebelum push pertama.

## ADR-002: Gin sebagai router
- Konteks: butuh server-rendered + session cookie + `templ.Handler`.
- Keputusan: Gin (bukan stdlib/Echo/Fiber).
- Alasan: middleware/session matang, tim umum memakainya, minim lock-in.

## ADR-003: GORM sqlite->postgres via env
- Konteks: dev SQLite, deploy Postgres Supabase tanpa ubah kode.
- Keputusan: satu `internal/db.Open()` switch dialek; larang SQL dialek-spesifik.
- Alasan: natur GORM; lihat `ARCHITECTURE.md:2`.

## ADR-004: templ untuk HTML- Konteks: butuh komponen type-safe dipetakan dari `DESIGN.md:569`.
- Keputusan: a-h/templ, `templ generate` wajib sebelum build.
- Alasan: cek compile-time, cocok dengan pola server-rendered + HTMX partial.

## ADR-005: Strategi test proyek besar- Konteks: butuh banyak unit test tanpa layanan luar.
- Keputusan: stdlib `testing` + table-driven; DB test pakai SQLite di `t.TempDir()`; `make test` sebagai gerbang (`go test ./...`).
- Alasan: tanpa dep tambahan, hermetis, pola sama dipakai ulang untuk sizematch/poin/kode nanti.

## ADR-006: Logika murni terpisah dari DB
- Konteks: charset kode, TTL 60 mnt, tabel poin fully specified di DESIGN.
- Keputusan: `internal/orders` (Generate/ValidFormat/ExpiredAt) + `internal/points` (EarnFor/RedeemCost/Apply/CanRedeem) murni tanpa GORM, unit-test penuh.
- Alasan: aturan bisnis terverifikasi tanpa DB; service tinggal panggil.

## ADR-007: Toleransi size-match 10% + skor aturan
- Konteks: rumus % dan metode skor TBD di DESIGN; user putuskan: margin error ~10% agar pilihan banyak; skor dari kelengkapan.
- Keputusan: per dimensi `1 - |item-user|/(0.10*user)`, rata-rata, bulat; Complete wajib Dada+Torso; Tinggi tak dinilai. Skor kondisi = `100-10*cacat` (lantai 40), perlu 6 foto, method "aturan".
- Alasan: deterministik, testable, tanpa klaim AI; ambang bisa di-tune satu konstanta.

## ADR-008: UI interim tanpa build Tailwind
- Konteks: biner tailwindcss belum ada; halaman harus sudah bisa diverifikasi.
- Keputusan: `app.css` tulis tangan dari token DESIGN + `input.css` sebagai sumber `make css` nanti; htmx/alpine via CDN (vendor lokal saat prod).
- Alasan: unblock render/UX review; token tetap satu sumber.

## ADR-009: Sesi cookie HMAC tanpa dep tambahan
- Konteks: butuh session filler 30 hari, belum OTP/OAuth.
- Keputusan: cookie `userID.expiry.sig` HMAC-SHA256 secret; verifikasi + lookup di middleware.
- Alasan: stdlib saja, murni/testable; ganti provider nanti tanpa ubah handler.

## ADR-010: Aturan transaksi
- Konteks: cegah double-sell tanpa kolom hold; poin hanya pasca-done.
- Keputusan: guard satu pending per produk; cart dibersihkan saat checkout; confirm wajib kode lawan; done (2 sisi) → produk sold + ledger +10 pembeli; p2p selesai manual tanpa poin.
- Alasan: tanpa status tambahan; sengketa/expired lanjutan tetap TBD produk.

## ADR-011: Aturan sosial + geo + redeem
- Konteks: ulasan/chat/lapor/jarak/penukaran butuh guard backend.
- Keputusan: ulasan = order done + reviewer buyer + 1/order; rating kosong → "Belum ada rating". Thread chat dipakai ulang per produk+order; pengirim harus anggota. Jarak Haversine + format koma ID. Redeem atomik (cek saldo + catat + ledger negatif).
- Alasan: eligibility di server, bukan UI; saldo selalu dari ledger.

## ADR-012: Aturan lapis template anti-siklus
- Konteks: dua siklus impor Go tertangkap vet (pages↔catalog, components↔pages, db↔catalog via test).
- Keputusan: data bersama di `models` (Categories); helper format di `components`; seed demo di paket `seed` sendiri; handler boleh impor service lain selama asiklik.
- Alasan: arah impor tunggal pages→components, handler→service; siklus = sinyal salah tempat.

## ADR-013: Upload via Storage interface + sniffing
- Konteks: multipart dari browser tak bisa dipercaya ekstensinya.
- Keputusan: `storage.Local` validasi via `http.DetectContentType` + batas 5MB + nama acak; DB simpan URL string (ganti Supabase tanpa migrasi).
- Alasan: smoke buktikan SVG ditolak, PNG lolos, publish 302 ke produk aktif.

## ADR-014: Gate lokasi + chat per order
- Konteks: DESIGN wajibkan lokasi sebelum bayar; chat butuh konteks transaksi.
- Keputusan: checkout diblokir bila lat+lng nol (server-side, bukan cuma UI); thread chat dibuat per order di halaman transaksi/konfirmasi; pesan polling HTMX 5 dtk via partial.
- Alasan: smoke buktikan kode dua arah → sukses +10 poin; tanpa WebSocket untuk mock.

## ADR-015: Lookup user terpusat + kaidah template
- Konteks: siklus auth↔users; sintaks `@{ }` ilegal di templ.
- Keputusan: FindByEmail/FindByID pindah ke `auth` (users boleh impor auth, bukan sebaliknya); teks `@handle` via helper `components.At`.
- Alasan: arah impor satu arah; `@` mentah + `{expr}` berdampingan ditolak parser templ.

## ADR-016: Admin minimum tanpa suspend
- Konteks: PRD minta moderasi; kebijakan suspend/banding TBD.
- Keputusan: `users.is_admin` + `RequireAdmin`; baca (ringkasan/transaksi) + resolve/reject laporan + takedown (bukan hapus) + audit tiap aksi. Tanpa suspend/blokir/hapus permanen.
- Alasan: smoke buktikan non-admin 403; peta statis + pin beda + jarak koma-ID.

## ADR-017: Rute akun di paket users
- Konteks: signup butuh users.Register, tapi auth↔users siklus.
- Keputusan: /signup + /lupa-sandi di users/handler; sesi via auth.NewSession; lupa-sandi selalu netral.
- Alasan: smoke buktikan daftar → preferensi; tak ada email bocor.

## ADR-018: CSS dibuild, audit-UI ditunda
- Konteks: `make css` butuh resolusi paket global; audit tercatat tapi tanpa halaman daftar.
- Keputusan: `input.css` = `@theme` + gaya komponen; build via NODE_PATH global (repo bebas node_modules). Daftar audit jadi follow-up, bukan MVP.
- Alasan: `make css` + render terverifikasi; satu sumber token tetap.
