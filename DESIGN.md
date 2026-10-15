# DESIGN.md — MEWALI

> **Sumber acuan utama untuk desain UI/UX dan implementasi antarmuka MEWALI.** Baca seluruh dokumen ini sebelum membuat atau mengubah UI. Jika ada konflik dengan default framework atau preferensi pribadi coding agent, ikuti dokumen ini. PDF desain terbaru menjadi acuan visual untuk layar yang sudah dirancang; PRD menjadi acuan cakupan produk. Jangan mengarang layar atau aturan bisnis yang belum ditentukan: tandai sebagai keputusan yang perlu dikonfirmasi.

## 1. Ringkasan produk

**MEWALI** adalah platform jual-beli pakaian thrift/preloved lokal yang membantu pakaian mendapat kehidupan kedua melalui transaksi yang lebih jelas dan terpercaya.

- **Tagline:** `NEW LIFE, SAME STYLE`
- **Visi:** Menjadi platform thrifting terpercaya di Indonesia agar pengguna dapat berbelanja dengan yakin, menemukan ukuran yang sesuai, dan memahami dampak positif pilihannya.
- **Misi:** Memudahkan penggunaan kembali pakaian melalui jual-beli yang praktis serta meningkatkan kepercayaan lewat informasi kondisi dan cacat yang transparan.
- **Pengguna utama:** pembeli, penjual, dan admin/pengelola platform.
- **Prinsip pengalaman:** terpercaya, hangat, praktis, transparan, inklusif, dan berorientasi pada penggunaan kembali.

### Produk dan permukaan antarmuka

1. **Aplikasi pengguna, mobile-first:** pembeli dan penjual menggunakan akun yang sama. Pengguna dapat membeli sekaligus menjual; jangan memisahkan keduanya menjadi aplikasi berbeda.
2. **Dashboard admin, web desktop:** antarmuka terpisah untuk pengelolaan platform, moderasi konten, dan pengawasan transaksi. Dashboard admin bukan sekadar tampilan mobile yang diperlebar.
3. **Desain responsif:** aplikasi pengguna mengutamakan lebar acuan 390 px, tetapi tetap harus berfungsi pada layar kecil dan besar. Dashboard admin mengutamakan desktop, dengan layout tablet yang tetap layak.

## 2. Prioritas sumber dan aturan desain

Jika ada perbedaan, gunakan urutan berikut:

1. **PRD** menentukan kebutuhan produk dan daftar fitur.
2. **PDF desain terbaru** menentukan tampilan layar yang telah dirancang.
3. **Dokumen ini** menyatukan keduanya dan menjelaskan pola yang belum terperinci.
4. Keputusan implementasi teknis tidak boleh menghapus kebutuhan produk atau mengubah desain yang telah disetujui tanpa alasan.

- Layar yang tampak di PDF harus dipertahankan secara visual sejauh memungkinkan.
- Fitur yang diminta PRD tetapi belum memiliki desain PDF harus mengikuti design system dokumen ini.
- Jangan menganggap contoh nama, produk, harga, rating, lokasi, kode transaksi, atau isi ulasan sebagai data produksi.
- Jangan membuat data statis seolah-olah data nyata. Gunakan data API, fixture/demo yang jelas, atau empty state.
- Jangan menambahkan branding atau elemen dari frame wireframe yang bukan bagian MEWALI.

## 3. Bahasa, penulisan, dan format konten

- Bahasa UI utama: **Bahasa Indonesia**.
- Sapaan: **“kamu”**, bukan “Anda”. Nada hangat, singkat, jelas, dan tidak menghakimi.
- Judul layar mengikuti gaya PDF: **HURUF KAPITAL**, kecuali judul yang memang berupa kalimat atau bagian konten.
- Gunakan ejaan yang benar dan konsisten: “di kemudian hari”, “Berikan kode”, “pembeli”, “penjual”, “Tambah katalog”. Jangan menyalin typo wireframe.
- Mata uang: `Rp65.000` pada detail dan ringkasan transaksi; harga kartu dapat diringkas menjadi `65K` atau `40K`.
- Desimal: gunakan koma, misalnya `4,8` dan `0,8 km`.
- Ukuran pakaian: tampilkan satuan `cm` dan jelaskan apakah angka adalah ukuran barang atau ukuran pengguna.
- Tanggal/waktu menggunakan format Indonesia yang mudah dibaca.
- Label tombol harus menyebut tindakan yang akan terjadi, misalnya “Simpan preferensi”, “Unggah foto”, “Bayar”, “Kirim ulasan”. Hindari tombol generik seperti “OK” jika tindakan bisa dijelaskan.
- Placeholder nama dan konten dalam PDF adalah contoh saja, bukan teks wajib produksi.

## 4. Identitas visual dan design tokens

### 4.1 Filosofi brand

- **Panah berputar:** sirkulasi dan penggunaan kembali.
- **Baju/gantungan baju:** pakaian yang siap dipajang dan dipindahtangankan.
- **Tumpukan pakaian:** pemilahan dan pengumpulan pakaian preloved.
- **Pakaian/totebag:** kepemilikan baru dan pengalaman belanja.
- **Warm antique gold:** nilai yang tetap berharga meskipun barang pernah dimiliki.
- **Sea green/forest green:** alam dan keberlanjutan.
- **Slate/denim blue, warm taupe/khaki, muted maroon/burgundy, terracotta/soft coral:** aksen terinspirasi warna pakaian dan fashion.

Logo MEWALI adalah wordmark khusus; gunakan aset SVG/vektor resmi tim. Jangan mengetik ulang logo sebagai teks atau membuat ulang bentuk huruf W secara manual.

### 4.2 Token warna

Gunakan token terpusat; jangan menulis hex langsung di komponen UI. Warna di bawah mempertahankan nuansa PDF/desain lama dan menjadi baseline. Jika tim desain menyediakan token final dari Figma, perbarui nilai token di satu tempat tanpa mengganti nama token di semua komponen.

```css
:root {
  /* Base */
  --bg: #F5F5F5;
  --surface: #EFEFEF;
  --surface-raised: #FFFFFF;
  --surface-muted: #E7E8E5;
  --text: #222222;
  --text-muted: #8A8A8A;
  --border: #D9DEDA;

  /* Brand */
  --primary: #78A298;       /* sage/sea green */
  --primary-strong: #456F63; /* forest green untuk teks/kontras */
  --primary-ink: #FFFFFF;
  --brown: #6B4A3A;
  --gold: #F5C518;
  --gold-warm: #C6A04A;     /* warm antique gold */

  /* Aksen kategori / status */
  --denim: #607D9A;
  --taupe: #B5A38A;
  --burgundy: #87515C;
  --terracotta: #C97860;
  --success: #3F7A58;
  --warning: #A66A18;
  --danger: #B5432F;
  --danger-soft: #A97070;
  --info: #607D9A;

  /* Dekorasi dan chat */
  --blob-edge: #E3E9E7;
  --blob-center: #F4F2EE;
  --bubble-in: #D9D9D9;
  --bubble-out: #D6F3EC;

  /* Radius */
  --r-pill: 999px;
  --r-card: 24px;
  --r-input: 20px;
  --r-chip: 14px;
  --r-modal: 28px;

  /* Spacing */
  --page-x: 18px;
  --gap: 10px;
  --section-gap: 24px;

  /* Controls */
  --h-btn: 56px;
  --h-input: 56px;
  --tap-min: 44px;
}
```

**Penggunaan warna:**

- `--primary` untuk aksi utama, ikon aktif, progress, dan elemen brand.
- `--primary-strong` untuk teks atau elemen hijau yang membutuhkan kontras lebih kuat.
- `--gold` untuk bintang/poin dan aksen terbatas; jangan jadikan warna latar utama.
- `--brown` untuk aksen hangat dan badge poin.
- `--danger` hanya untuk cacat, pelaporan, validasi gagal, atau tindakan berisiko.
- Warna kategori adalah aksen pendukung, bukan alasan untuk mewarnai setiap kartu secara acak.
- Pastikan rasio kontras teks dan latar memadai. Jangan mengasumsikan teks putih di atas sage green otomatis lolos; gunakan `--primary-strong` atau warna teks gelap bila perlu.

Jika memakai Tailwind CSS v4, petakan token ke `@theme` dan gunakan kelas seperti `bg-primary`, `text-primary-strong`, `bg-surface`. Jangan menyebar nilai hex di template.

## 5. Tipografi, ikon, ilustrasi, dan dekorasi

- **Display/heading:** font lebar dan tegas seperti Syne atau Krona One, jika aset/font telah tersedia. Gunakan kapital, bobot bold, dan letter spacing ringan (`.04em`) pada judul layar.
- **Isi, input, tombol:** Plus Jakarta Sans atau sans-serif geometris yang setara; ukuran isi umumnya 14–16 px.
- **Teks sekunder:** 12–13 px, warna muted, tetapi jangan gunakan ukuran kecil untuk informasi penting, status transaksi, cacat, atau validasi.
- Nama produk pada kartu menggunakan display bold; deskripsi tetap mudah dibaca dan tidak terlalu padat.
- Gunakan ikon garis dengan gaya konsisten. Hindari mencampur ikon outline dan ikon filled tanpa alasan.
- Ikon tanpa teks wajib memiliki label aksesibilitas.
- Blob dekoratif berupa bunga/bintang enam kelopak membulat dengan gradien netral; gunakan terutama pada onboarding, autentikasi, preferensi, dan konfirmasi.
- Bintang emas dengan glow hanya digunakan pada banner atau elemen poin yang relevan.
- Dekorasi tidak boleh mengurangi keterbacaan, menutupi kontrol, atau mengganggu navigasi.
- Gunakan ilustrasi/aset dari tim jika tersedia; jangan menggunakan logo/ilustrasi generik sebagai pengganti aset resmi tanpa kebutuhan.

## 6. Layout, responsivitas, dan komponen dasar

### 6.1 Layout aplikasi pengguna

- Mobile-first, referensi 390 px; gunakan padding horizontal sekitar 18 px.
- Beri jarak antarbagian yang konsisten dan hindari tampilan terlalu rapat.
- Halaman dengan bottom navigation perlu padding bawah sekitar 110 px agar konten tidak tertutup.
- Hormati safe area perangkat, terutama dock dan tombol aksi bawah.
- Gunakan grid dua kolom untuk katalog pada ponsel. Pada layar lebih besar, tambah kolom secara bertahap tanpa membuat kartu terlalu lebar.
- Elemen interaktif memiliki target sentuh minimal 44 × 44 px.
- Hindari horizontal overflow, teks terpotong tanpa affordance, dan tombol yang tidak terlihat saat keyboard terbuka.

### 6.2 Tombol

- Tombol utama berbentuk pill, tinggi 56 px, lebar penuh bila menjadi CTA utama.
- Tombol utama menggunakan `--primary` dengan warna teks yang lolos kontras.
- Tombol sekunder dapat memakai `--surface` dan teks `--primary-strong`.
- Tombol destructive menggunakan `--danger` atau outline yang jelas; jangan memakai merah untuk aksi biasa.
- Pasangan “Kembali” + aksi lanjut memakai dua kolom, kira-kira `1fr 1.4fr`, dengan gap 10 px.
- Bedakan state default, hover/focus, pressed, disabled, loading, dan sukses.
- Tombol disabled harus terlihat nonaktif dan tidak menjalankan aksi.

### 6.3 Input dan form

- Input memakai latar `--surface`, tanpa border berat, radius `--r-input`, tinggi 56 px, padding horizontal sekitar 18 px.
- Ikon garis di kiri bila membantu identifikasi field; jangan menambahkan ikon dekoratif ke setiap field jika tidak berguna.
- Focus terlihat jelas, misalnya outline 2 px warna primary.
- Label tetap terlihat saat pengguna mengisi field; placeholder bukan satu-satunya label.
- Tampilkan validasi dekat field yang bermasalah dengan bahasa yang menjelaskan cara memperbaiki.
- Field ukuran menerima angka dan menunjukkan satuan cm. Validasi nilai kosong, negatif, dan format tidak masuk akal.
- Sediakan state loading, error, dan success untuk submit form.

### 6.4 Chip kategori

- Bentuk persegi membulat dengan latar `--surface`, ikon garis dan label.
- Chip terpilih memiliki pembeda yang jelas menggunakan warna brand, bukan hanya perubahan warna yang samar.
- Daftar kategori awal: **Kaos, Kemeja, Jeans, Gaun, Jaket**, serta kategori tambahan yang tersedia dari data.
- Daftar horizontal boleh scroll, tetapi harus tetap dapat digunakan dengan sentuhan dan keyboard.

### 6.5 Kartu produk

- Foto menjadi elemen visual utama; rasio kartu 3:4 dan `object-fit: cover`.
- Radius `--r-card` (24 px), overflow tersembunyi.
- Harga berupa badge pill terang di kanan atas.
- Bagian bawah menggunakan gradien gelap lembut agar nama dan penjual putih tetap terbaca. Jangan mengandalkan blur saja untuk kontras.
- Nama produk dan penjual ditempatkan di kiri bawah; batasi baris dengan ellipsis bila diperlukan.
- Tampilkan status penting (terjual, ditahan, draft) dengan badge yang jelas pada konteks kelola produk, bukan secara otomatis pada semua kartu katalog publik.
- Informasi kecocokan ukuran atau kondisi dapat ditambahkan bila relevan, tanpa menjejalkan terlalu banyak badge.
- Gunakan `loading="lazy"` untuk gambar di luar viewport dan teks alternatif yang menjelaskan produk.

### 6.6 Bottom navigation / dock

- Kapsul melayang di bawah konten, latar `--surface`, radius pill, memperhitungkan safe area.
- Item aktif berupa pill `--primary` yang menampilkan ikon dan label; item tidak aktif dapat menampilkan ikon saja.
- Navigasi utama: **Beranda · Peta · Profil**, dengan tombol **Cari** terpisah bila mengikuti layout PDF.
- Akses **Jual** harus mudah ditemukan dari profil/barang dijual dan/atau CTA yang konsisten. Jangan mengubah tiga item navigasi utama menjadi penuh sesak; gunakan FAB atau entry point yang jelas.
- Pastikan state aktif sesuai route dan seluruh ikon memiliki label aksesibilitas.

### 6.7 Pola status umum

Setiap fitur yang mengambil data harus memiliki state berikut jika relevan:

- **Loading:** skeleton atau indikator sederhana, bukan layout kosong yang tampak rusak.
- **Empty:** pesan yang menjelaskan kondisi dan CTA berikutnya, misalnya “Belum ada barang yang dijual” + “Tambah katalog”.
- **Error:** pesan yang bisa dipahami dan aksi coba lagi jika sesuai.
- **Success:** konfirmasi tindakan selesai.
- **Disabled / unavailable:** jelaskan alasan aksi belum dapat dilakukan.
- **Permission denied:** jelaskan kenapa izin lokasi/kamera dibutuhkan dan berikan alternatif bila mungkin.

## 7. Navigasi dan arsitektur layar

### Aplikasi pengguna

1. Onboarding
2. Sign up, login, lupa sandi
3. Preferensi ukuran dan preferensi gaya
4. Beranda dan eksplorasi
5. Pencarian, filter, dan hasil katalog
6. Detail produk
7. Keranjang dan checkout
8. Status transaksi dan konfirmasi dua arah
9. Chat antar-pengguna
10. Profil sendiri dan profil pengguna lain
11. Barang dijual / pengelolaan katalog penjual
12. Tambah dan edit katalog
13. Poin & Dampak serta penukaran reward
14. Peta penjual dan titik donasi
15. Rating, ulasan, dan pelaporan
16. Pengaturan akun dan preferensi

Pengguna yang sama dapat bertindak sebagai pembeli dan penjual. Pertahankan konteks saat berpindah dari profil ke katalog atau transaksi; jangan memaksa pengguna membuat akun kedua.

### Dashboard admin (web terpisah)

- Login admin dan pengelolaan sesi.
- Ringkasan operasional platform.
- Moderasi produk/listing dan konten yang dilaporkan.
- Pengawasan transaksi dan status sengketa/masalah yang tersedia pada sistem.
- Pengelolaan pengguna dan tindakan moderasi yang diizinkan.
- Pengelolaan kategori, titik donasi, opsi reward, atau konfigurasi platform hanya bila memang termasuk cakupan implementasi.
- Ringkasan dampak/indikator platform bila data tersedia.
- Riwayat tindakan admin/audit log untuk aksi moderasi penting.

**Catatan:** PRD menyebut pengelolaan platform, moderasi konten, dan pengawasan transaksi, tetapi belum menetapkan seluruh kolom, izin, atau prosedur admin. Buat layar minimum yang mendukung kebutuhan tersebut; tandai aturan operasional yang belum diputuskan, jangan mengarang kebijakan.

## 8. Spesifikasi layar aplikasi pengguna

### 8.1 Onboarding

Empat tampilan/step seperti pada PDF:

1. **Splash:** wordmark MEWALI besar di tengah blob, tagline `NEW LIFE, SAME STYLE`, dan teks “Powered by GROK” hanya jika itu memang bagian aset/brand final yang disetujui.
2. **Selamat Datang:** ilustrasi daur ulang dengan hanger/kaos dan panah hijau melingkar; judul “SELAMAT DATANG di MEWALI”, teks “Kami senang kamu di sini!”, CTA “Selanjutnya”.
3. **Apa Itu MEWALI?:** visual tumpukan pakaian dan ringkasan produk tentang jual-beli thrift yang aman, nyaman, dan mendukung pengurangan limbah tekstil. Tombol “Kembali” dan “Selanjutnya”.
4. **Dampak Kita:** ilustrasi/foto penyerahan pakaian dan kotak donasi; pesan bahwa jual-beli dan donasi memberi pakaian kesempatan kedua. Tombol navigasi onboarding.

Gunakan indikator titik; titik aktif memanjang menjadi pill cokelat. Onboarding tidak boleh menghalangi pengguna kembali ke step sebelumnya.

### 8.2 Autentikasi

**Sign up**

- Judul `SIGN UP`, subjudul “Jadi bagian dari komunitas kami!”.
- Field Nama, Email, Sandi, Konfirmasi Sandi.
- Opsi “Lanjutkan dengan Google/Apple” hanya jika autentikasi provider benar-benar tersedia.
- Link “Sudah punya akun? Login”, syarat dan ketentuan, CTA “Buat Akun”.

**Login**

- Judul `LOGIN`, subjudul “Selamat datang kembali”.
- Field Email dan Sandi; link “Lupa Sandi?”.
- Opsi Google/Apple bila terintegrasi, link ke Sign up, syarat dan ketentuan, CTA “Login”.

**Lupa sandi**

- Tombol kembali bulat, ikon kunci, judul `LUPA SANDI?`, instruksi pengiriman tautan reset, field Email, CTA “Kirim”.
- Setelah submit, tampilkan konfirmasi netral yang tidak membocorkan apakah email terdaftar.

**Preferensi ukuran dan gaya**

- Judul `PREFERENSI` dan penjelasan singkat.
- Field ukuran tubuh: Lebar Dada, Lingkar Pinggang, Lingkar Pinggul, Tinggi Badan, Panjang Torso, dalam cm.
- Sediakan preferensi gaya/kategori jika data dan desain produk mendukungnya.
- CTA yang tepat adalah “Simpan” atau “Simpan Preferensi”; label “Login” pada wireframe lama dianggap perlu diperbaiki, bukan perilaku yang harus ditiru.
- Pengguna dapat melewati atau memperbarui preferensi bila kebijakan produk mengizinkan; jangan membuat preferensi tubuh wajib tanpa keputusan produk.

### 8.3 Beranda

- Header: wordmark MEWALI di kiri, tombol keranjang di kanan.
- Carousel banner dengan kartu permukaan membulat; kartu berikutnya boleh sedikit terlihat di tepi.
- Chip kategori horizontal.
- Seksi **Barang Baru** dan **Rekomendasi**, masing-masing memakai grid dua kolom dan tautan “Lihat Lebih”.
- Kartu produk menampilkan foto, harga ringkas, nama, dan penjual.
- Banner yang dapat digunakan:
  - “Terdekat dengan dirimu” dengan CTA “Beli”.
  - “Raih Poin / Bantu Sesama dan Sekitar” dengan CTA “Lihat Poin”.
  - “Jangan Akhiri Ceritanya di Lemari” dengan CTA “Mulai Jual”.
- Banner terdekat tidak boleh mengklaim lokasi atau jarak sebelum data lokasi tersedia. Jika izin lokasi belum ada, berikan penjelasan dan pilihan mengatur lokasi.

### 8.4 Pencarian dan eksplorasi

- Judul `PENCARIAN`.
- Search field dengan tombol kembali, ikon pencarian, dan tombol hapus.
- Grid produk dua kolom dan chip kategori.
- Kebutuhan PRD mencakup mesin pencarian pakaian; sediakan filter yang relevan berdasarkan kategori, rentang harga, ukuran/kecocokan, kondisi, lokasi/jarak, dan urutan hasil bila filter tersebut didukung backend.
- Tampilkan jumlah/hasil yang sesuai bila tersedia, state loading, hasil kosong, dan kegagalan pencarian.
- Jangan mengklaim hasil “paling cocok” jika kriteria ranking belum diimplementasikan.

### 8.5 Detail produk

Urutan konten yang disarankan:

1. Galeri foto produk dengan tombol kembali dan keranjang melayang, indikator slide, serta kontrol yang mudah digunakan.
2. Nama produk, harga penuh (`Rp65.000`), waktu foto/unggahan bila tersedia.
3. Ringkasan penjual: avatar/ikon toko, nama, rating, jarak bila lokasi tersedia, dan tautan ke profil.
4. Skor/keterangan kondisi barang yang jelas.
5. Kecocokan ukuran pembeli dengan barang: persentase dan perbandingan ukuran dalam cm.
6. Deskripsi produk.
7. Foto dengan penanda cacat dan daftar cacat yang terhubung ke penanda pada foto.
8. CTA bawah: tambah ke keranjang dan “Beli sekarang”.

**Penandaan cacat:** nomor penanda pada foto harus cocok dengan daftar teks bernomor. Gunakan warna `--danger` dengan bentuk yang cukup kontras; penanda tidak boleh menutupi seluruh cacat. Setiap cacat memiliki lokasi/area dan deskripsi dari penjual. Jika belum ada cacat yang ditandai, jangan menyimpulkan barang sempurna; tampilkan “Penjual belum menandai cacat” atau teks yang sesuai dengan data.

**Kecocokan ukuran:** tampilkan persentase hanya jika dihitung oleh logika yang disepakati dan ukuran pembeli serta barang tersedia. Tampilkan bar visual dan perbandingan yang mudah dibaca, misalnya “Dada barang / kamu — 54 / 52 cm”. Jika data tidak lengkap, jelaskan bahwa kecocokan belum dapat dihitung dan arahkan pengguna melengkapi preferensi.

### 8.6 Keranjang

- Judul `KERANJANG`.
- Item menampilkan foto, nama produk, penjual, harga, checkbox/pemilihan item, dan badge poin yang akan diperoleh jika sesuai aturan.
- Ringkasan total jelas dan CTA “Beli”.
- Sediakan empty state dengan CTA kembali ke katalog.
- Perubahan jumlah/item harus memperbarui subtotal dan total; jangan tampilkan angka statis yang tidak mengikuti data.

### 8.7 Checkout / Beli

- Judul `BELI`.
- Tampilkan alamat/lokasi pengambilan atau transaksi. Jika lokasi belum disetel, tampilkan peringatan yang jelas dengan tautan untuk mengatur lokasi.
- Ringkasan produk, jumlah item, total produk, subtotal pesanan, dan total akhir.
- Metode pembayaran mengikuti PRD/desain yang benar-benar tersedia. PDF menampilkan **Cash** dan **Transfer ke Penjual**; jangan menampilkan payment gateway/provider yang belum terintegrasi seolah-olah aktif.
- CTA “Bayar” atau CTA sesuai metode pembayaran dan status transaksi.
- Jelaskan biaya tambahan sebelum konfirmasi bila ada; jangan menyembunyikan biaya.
- **Aturan desain/perilaku:** pengguna tidak boleh melanjutkan pembelian yang mensyaratkan lokasi sebelum lokasi valid tersedia. Bila transaksi online/tanpa lokasi nantinya didukung, aturan ini perlu dikonfirmasi dan tidak boleh diasumsikan.

### 8.8 Konfirmasi transaksi dua arah

PDF menunjukkan konfirmasi dengan kode untuk transaksi tatap muka. Pertahankan pola ini selama mekanisme tersebut menjadi aturan bisnis aktif.

**Tampilan pembeli:**

- Judul `KONFIRMASI`.
- Kode milik pembeli ditampilkan dengan jelas beserta instruksi untuk memberikannya kepada penjual jika barang sesuai.
- Instruksi memasukkan kode dari penjual setelah barang diterima.
- Timer batas waktu, input kode lawan, CTA “Konfirmasi”, dan CTA “Chat Penjual”.

**Tampilan penjual:**

- Pola sama, dengan instruksi memberikan kode kepada pembeli dan memasukkan kode pembeli setelah penyerahan dinyatakan aman.
- CTA “Chat Pembeli”.

**Aturan interaksi:**

- Masing-masing pihak memiliki kode 5 karakter dan memasukkan kode milik pihak lain.
- Gunakan karakter yang tidak mudah tertukar; hindari pasangan `O/0` dan `I/1`.
- Batas waktu pada dokumen lama adalah 60 menit; tampilan dapat menunjukkan `59m` setelah hitung mundur berjalan.
- Kode tidak boleh hardcode atau dapat dipakai ulang.
- Tombol konfirmasi harus menjelaskan kesalahan kode dan status menunggu secara jelas.
- Transaksi selesai hanya setelah kedua pihak berhasil mengonfirmasi sesuai aturan backend.
- Poin transaksi baru diberikan setelah transaksi terkonfirmasi di kedua sisi.
- Jika timer habis, jaringan gagal, atau salah satu pihak membatalkan, tampilkan status dan tindakan berikutnya berdasarkan kebijakan produk; kebijakan penyelesaian kasus ini belum ditetapkan dalam PRD.

### 8.9 Halaman sukses transaksi

- Visual centang pada blob.
- Judul “Pembelian Selesai”, ringkasan poin yang diperoleh, dan pesan singkat tentang memperpanjang siklus hidup pakaian.
- CTA “Kembali ke Beranda” dan, bila sesuai, akses ke detail transaksi.
- Jangan menampilkan poin sebagai sudah didapat sebelum backend mengonfirmasi transaksi dan pemberian poin.

### 8.10 Chat antar-pengguna

- Tampilkan identitas lawan bicara, konteks produk/transaksi bila ada, daftar pesan, waktu kirim, serta input pesan.
- Pesan masuk memakai `--bubble-in`; pesan keluar memakai `--bubble-out`.
- Gelembung masuk/keluar dibedakan posisi dan bentuknya, bukan warna saja.
- Footer memiliki input pill “Tulis pesan...” dan tombol tambah/attachment hanya jika fitur lampiran benar-benar tersedia.
- Bedakan status loading, gagal mengirim, terkirim, dan pesan baru.
- Jangan menampilkan kode transaksi sensitif secara otomatis pada preview notifikasi atau konteks yang tidak diperlukan.
- Berikan entry point “Chat Penjual/Pembeli” dari transaksi dan produk sesuai konteks.

### 8.11 Profil sendiri

- Judul `PROFIL`.
- Avatar bundar dengan aksen primary, nama, `@handle`, jumlah produk, dan rating jika ada.
- Seksi **Barang Dijual** dengan tautan “Lihat Lebih”; empty state “Belum ada barang yang dijual” dan CTA “Tambah katalog”.
- Seksi **Poin & Dampak** menampilkan total poin dan level pengguna, misalnya “Pemula Hijau”, jika level tersebut tersedia.
- Seksi **Ulasan** dengan empty state yang ramah.
- Sediakan pengaturan akun dan preferensi di entry point yang mudah ditemukan.
- Pengguna dapat berpindah ke mode/fitur penjual tanpa akun terpisah.

### 8.12 Profil pengguna lain

- Informasi publik yang relevan: avatar, nama/handle, jumlah produk, rating, produk dijual, dan ulasan.
- Tombol “Laporkan” menggunakan gaya outline `--danger-soft` atau `--danger`.
- Tombol “Tambahkan Ulasan” hanya aktif jika pengguna memenuhi syarat transaksi/ulasan sesuai backend.
- Ulasan menampilkan penulis, skor bintang, dan isi.
- Jangan mengekspos data privat seperti email, alamat rinci, atau preferensi tubuh.
- Modal pelaporan meminta alasan dan keterangan secukupnya; jelaskan bahwa laporan akan ditinjau, tanpa menjanjikan hasil tertentu.

### 8.13 Barang dijual / daftar penjualan

- Judul `BARANG DIJUAL`.
- Header penjual, jumlah produk, rating, daftar produk, dan empty state.
- FAB `+` atau CTA “Tambah katalog” di posisi yang mudah dijangkau.
- Pada daftar milik sendiri, tampilkan status draft/aktif/terjual/ditinjau bila status tersebut tersedia. Sediakan tindakan edit/nonaktifkan/hapus sesuai izin.
- Pastikan pengguna tidak salah mengira draft atau listing yang belum lolos pemeriksaan sebagai barang aktif di marketplace.

### 8.14 Tambah dan edit katalog

Form harus mendukung:

- Gambar utama dan nama barang.
- Kategori: Kaos, Kemeja, Jeans, Gaun, Jaket, serta pilihan yang disediakan sistem.
- Foto wajib sesuai spesifikasi desain awal: **Tampak Depan, Tampak Belakang, Kerah, Lengan, Tag Samping, Tag Atas**. Jika tim produk mengubah kebutuhan foto, ubah konfigurasi/label secara konsisten.
- Ukuran barang: Lebar Dada, Lingkar Pinggang, Lingkar Pinggul, Tinggi Badan, Panjang Torso (cm), dengan definisi ukuran yang jelas untuk penjual.
- Harga, deskripsi, kondisi, dan lokasi bila diwajibkan oleh alur listing.
- Penandaan cacat pada foto: penjual memilih area pada gambar dan menambahkan deskripsi cacat; setiap penanda bernomor sesuai daftar cacat.
- CTA “Simpan Draf” bila draft didukung dan CTA untuk memeriksa/menerbitkan listing.
- Indikator “Cek Rating Kondisi” dengan progress dan hasil skor setelah foto yang diperlukan lengkap.
- Tip: “Skor akan muncul ketika kamu menambahkan semua foto dan melakukan cek rating kondisi, sebelum akhirnya kamu bisa menjual barangmu.”

**Aturan listing:**

- Penjual belum dapat menerbitkan produk sampai foto wajib diunggah dan pemeriksaan kondisi menghasilkan skor, sesuai rancangan saat ini.
- Beri feedback per foto yang kurang, bukan hanya pesan gagal generik.
- Skor kondisi harus disertai penjelasan ringkas tentang artinya; jangan menampilkan angka tanpa konteks.
- Jangan membuat klaim bahwa skor otomatis menilai kondisi dengan AI jika metode tersebut belum dipastikan oleh tim teknis/produk.
- Bedakan kondisi barang dari cacat yang dilaporkan penjual; keduanya saling melengkapi, bukan menggantikan satu sama lain.
- Saat edit, pertahankan nilai form dan foto yang sudah tersimpan. Beri konfirmasi sebelum tindakan yang berpotensi menghilangkan perubahan.

### 8.15 Poin & Dampak

- Judul `POIN & DAMPAK`.
- Kartu level: jumlah poin, nama level, dan progress ke level berikutnya jika aturan level sudah tersedia.
- Bagian “Cara mendapatkan poin” dalam baris/kartu permukaan.
- Reward yang ditampilkan pada desain awal:
  - Beli barang thrift: `+10`
  - Jual dengan foto lengkap dan ukuran asli: `+15`
  - Tanpa retur (pas dengan ukuranmu): `+5`
  - Serahkan barang rusak ke rework loop: `+20`
  - Ajak teman yang bertransaksi: `+25`
- Opsi penukaran pada desain awal:
  - Gratis Ongkir: `100 poin`
  - Donasi 1 pakaian layak: `80 poin`
  - Tanam 1 pohon: `120 poin`
- Badge poin berupa pill cokelat dengan teks emas, misalnya `+15 Poin`.
- Setiap reward menunjukkan jumlah poin, manfaat, syarat utama bila ada, dan CTA “Tukar”.
- Jika poin tidak cukup, tampilkan jumlah yang kurang; jangan membuat aksi “Tukar” terlihat sukses.
- Setelah penukaran, tampilkan konfirmasi dan perubahan saldo hanya berdasarkan respons sistem.
- Jangan menampilkan level/ambang progress yang belum ditentukan sebagai fakta. Level “Pemula Hijau” adalah contoh level yang sudah terlihat di desain; level lainnya dan ambang level belum ditentukan.
- Aturan retur, kedaluwarsa poin, pembatalan penukaran, dan validasi hadiah perlu dikonfirmasi jika akan diterapkan.

### 8.16 Peta penjual dan titik donasi

- Judul `PETA`; gunakan peta sebagai permukaan utama dengan warna pucat/desaturasi agar penanda terlihat.
- Penanda produk dapat memakai kartu mini dengan segitiga primary di bawah, mengikuti desain PDF.
- Penanda donasi harus berbeda dari penanda produk dan memiliki legenda yang jelas.
- Peta mendukung pencarian penjual/produk terdekat dan titik donasi terdekat, sesuai izin serta data lokasi.
- Saat penanda dipilih, tampilkan kartu ringkas: nama produk/titik donasi, informasi relevan, jarak jika tersedia, dan tindakan lanjut.
- Sediakan loading, peta gagal dimuat, lokasi belum diizinkan, hasil kosong, dan fallback daftar bila memungkinkan.
- Jangan meminta izin lokasi tanpa konteks. Jelaskan manfaatnya dan sediakan opsi memasukkan lokasi secara manual jika produk mendukungnya.
- Jangan tampilkan alamat rumah penjual secara presisi kepada publik kecuali memang merupakan lokasi publik yang secara sadar dibagikan penjual.
- Detail provider peta, geocoding, format koordinat, serta mekanisme privasi lokasi belum ditetapkan oleh PRD; jangan mengunci vendor dalam desain.

### 8.17 Rating dan ulasan

- Gunakan lima bintang outline warna primary; bintang terisi ketika dipilih.
- Modal/sheet memiliki judul, tombol kembali/tutup, area teks ulasan, rating, dan CTA “Kirim”.
- Tampilkan syarat ulasan jika hanya transaksi tertentu yang memenuhi syarat.
- Ulasan yang sudah terkirim memiliki state sukses; kegagalan submit mempertahankan teks pengguna.
- Hindari menampilkan rating sebagai `0,0` jika belum ada ulasan. Gunakan “Belum ada rating”.
- Aturan edit/hapus ulasan dan batas jumlah ulasan perlu dikonfirmasi bila diperlukan.

### 8.18 Laporkan pengguna atau produk

- Gunakan modal sheet dengan latar belakang diburamkan, tombol tutup/kembali, judul, pilihan alasan, field keterangan opsional bila sesuai, dan CTA “Kirim Laporan”.
- Jelaskan bahwa laporan dikirim untuk ditinjau.
- Setelah submit, tampilkan konfirmasi tanpa membuka detail laporan kepada pihak yang dilaporkan.
- Bedakan tindakan “Laporkan” dari tindakan destructive terhadap akun milik sendiri.

### 8.19 Pengaturan akun

- Sediakan akses untuk memperbarui profil, preferensi ukuran/gaya, lokasi, dan sesi akun.
- Konfirmasi sebelum logout atau tindakan yang dapat mengganggu sesi/transaksi berjalan.
- Jangan menganggap fitur yang belum ada di PRD sebagai wajib. Tambahkan pengaturan privasi/keamanan yang diperlukan untuk implementasi, tetapi tandai keputusan produk yang memengaruhi perilaku.

## 9. Dashboard web admin

### 9.1 Prinsip visual

- Web desktop dengan hierarki informasi yang padat namun teratur; jangan menyalin kartu produk mobile sebagai layout utama dashboard.
- Gunakan token warna dan tipografi yang sama agar tetap satu brand, tetapi prioritaskan keterbacaan data, tabel, filter, dan status.
- Gunakan sidebar navigasi atau pola navigasi desktop setara, header halaman, area konten, dan breadcrumb bila membantu.
- Tabel harus mendukung pemindaian cepat: header jelas, kolom terurut secara visual, status dalam badge, dan aksi pada posisi konsisten.
- Pada layar sempit, tabel dapat berubah menjadi kartu atau memiliki scroll horizontal yang disengaja.

### 9.2 Area dashboard minimum

**Ringkasan**

- Ringkasan jumlah pengguna, listing, transaksi, laporan, dan dampak hanya jika data metrik tersedia.
- Gunakan grafik hanya jika membantu melihat tren; label, periode, dan satuan harus jelas.

**Moderasi listing/konten**

- Daftar produk yang dilaporkan atau perlu ditinjau.
- Informasi inti: foto, nama produk, penjual, status, alasan laporan, waktu, dan tindakan yang diizinkan.
- Detail listing memperlihatkan foto, deskripsi, ukuran, skor kondisi, dan penanda cacat.
- Aksi moderasi harus memiliki konfirmasi dan alasan/jejak audit bila relevan.

**Pengawasan transaksi**

- Daftar transaksi dengan ID, pihak terkait, nilai, metode pembayaran, status, dan waktu sesuai data yang diizinkan.
- Filter berdasarkan status/waktu dan akses detail untuk memeriksa progres transaksi.
- Jangan menampilkan data sensitif yang tidak diperlukan untuk tugas admin.

**Pengguna dan laporan**

- Cari pengguna, lihat status akun serta ringkasan laporan yang relevan, dan lakukan tindakan hanya sesuai hak akses.
- Tindakan seperti suspend, blokir, atau penghapusan permanen tidak boleh ditambahkan sebagai kemampuan aktif tanpa keputusan kebijakan produk.

**Audit dan akses**

- Aksi admin penting sebaiknya tercatat: siapa melakukan apa dan kapan.
- UI harus membedakan akses read-only dan aksi yang mengubah data.
- Detail role/permission dan prosedur banding belum didefinisikan PRD; jangan mengarang kebijakan final.

## 10. Pola interaksi dan aturan bisnis yang harus tercermin di UI

1. **Pembeli sekaligus penjual:** satu akun mendukung dua peran; entry point jual dan beli mudah ditemukan.
2. **Kecocokan ukuran:** hitung berdasarkan preferensi pengguna dan ukuran produk sesuai logika yang disepakati. Jangan hardcode persentase contoh.
3. **Skor kondisi:** listing tidak dapat dijual sebelum foto wajib lengkap dan pengecekan kondisi selesai, sesuai desain saat ini.
4. **Penandaan cacat:** penanda foto dan daftar deskripsi selalu saling cocok.
5. **Lokasi checkout:** cegah kelanjutan alur yang membutuhkan lokasi jika lokasi belum valid; jelaskan langkah memperbaikinya.
6. **Konfirmasi dua arah:** kode 5 karakter unik untuk tiap pihak, masing-masing memasukkan kode pihak lain, dengan batas waktu 60 menit sesuai desain saat ini.
7. **Poin:** berikan poin hanya setelah kondisi transaksi yang disyaratkan terkonfirmasi; jangan mengandalkan tampilan frontend sebagai sumber kebenaran saldo.
8. **Ulasan:** jangan izinkan aksi yang melanggar aturan eligibility backend; jika syarat belum diketahui, minta keputusan produk.
9. **Reward:** periksa saldo dan syarat penukaran sebelum menyelesaikan aksi.
10. **Moderasi:** laporan harus memiliki status yang bisa dipahami pelapor tanpa membocorkan proses internal yang sensitif.
11. **Aksi asinkron:** semua tombol yang memicu request harus memberi feedback loading, success, dan error.
12. **Data dinamis:** harga, rating, jarak, jumlah produk, poin, status transaksi, dan kondisi harus berasal dari data aktual, bukan nilai contoh desain.

## 11. Aksesibilitas, keamanan, dan kepercayaan

- Target sentuh minimal 44 px; kontrol penting tidak boleh hanya dapat digunakan dengan hover.
- Pastikan navigasi keyboard dan focus ring terlihat, termasuk modal, sheet, dropdown, dan input kode.
- Semua gambar produk memiliki alt text yang berguna; gambar dekoratif memakai alt kosong.
- Jangan gunakan warna sebagai satu-satunya penanda status; sertakan label/ikon/teks.
- Error form harus terkait dengan field dan diumumkan secara aksesibel.
- Modal harus memiliki fokus yang terkelola, dapat ditutup dengan cara yang sesuai, dan tidak membiarkan interaksi latar saat terbuka.
- Jangan mengekspos email, kode transaksi, alamat privat, atau informasi sensitif di halaman publik.
- Jangan menampilkan klaim “terverifikasi”, “aman”, “asli”, atau skor otomatis sebagai jaminan mutlak jika sistem tidak mendukungnya.
- Informasi kondisi barang harus faktual dan transparan; hindari bahasa yang menyalahkan penjual atau pembeli.
- Konfirmasi aksi berisiko, termasuk menghapus listing atau melakukan tindakan moderasi.

## 12. Komponen yang sebaiknya reusable

Implementasikan sebagai komponen bersama sesuai arsitektur proyek, bukan menduplikasi markup dengan gaya berbeda:

- `Button` (primary, secondary, ghost, danger, loading, disabled)
- `IconButton`
- `Field`, `Textarea`, `Select`, `FormError`
- `CategoryChip`
- `ProductCard`, `ProductMiniCard`, `PriceBadge`
- `SellerSummary`, `RatingSummary`, `ReviewCard`
- `Dock` / bottom navigation
- `PointsBadge`, `PointsLevelCard`, `RewardCard`
- `MatchBar`, `MeasurementComparison`
- `ConditionScore`, `DefectMarker`, `DefectList`
- `PhotoUploadSlot`, `PhotoChecklist`, `ImageGallery`
- `ChatBubble`, `ChatComposer`
- `TransactionStatus`, `ConfirmationCodeInput`, `Countdown`
- `Modal` / `Sheet`, `ReportForm`, `ReviewForm`
- `EmptyState`, `LoadingState`, `ErrorState`, `PermissionState`
- Komponen dashboard: `AdminSidebar`, `PageHeader`, `DataTable`, `StatusBadge`, `FilterBar`, `ConfirmDialog`

Komponen reusable harus memiliki API yang jelas dan state visual yang konsisten. Hindari membuat komponen abstrak berlebihan untuk elemen yang hanya dipakai sekali dan tidak memiliki perilaku/visual bersama.

## 13. Pedoman implementasi frontend

Stack yang disebut pada desain awal: **HTML server-rendered dengan templ + Tailwind CSS + HTMX**. Gunakan stack yang telah dipilih proyek; jangan mengganti framework hanya karena preferensi coding agent.

- Petakan token desain ke tema Tailwind, jangan hardcode warna/radius/spasi dalam setiap template.
- Interaksi navigasi dapat memakai `hx-boost` bila sesuai arsitektur aplikasi.
- HTMX cocok untuk submit form, pemuatan partial, pencarian, filter, modal/sheet, dan pembaruan bagian halaman.
- JavaScript kecil dapat digunakan untuk carousel, penempatan marker pada foto, input kode, timer visual, dan interaksi yang memang membutuhkan state client-side.
- Interaksi peta, chat real-time, dan manipulasi gambar kompleks mungkin membutuhkan library/komponen khusus. Pilih solusi berdasarkan kebutuhan dan batas teknis, bukan memaksakan semua interaksi ke HTMX.
- Jangan menulis ulang seluruh halaman melalui JavaScript jika pola server-rendered sudah mencukupi.
- Pastikan tombol, form, navigasi, filter, upload foto, peta, dan modal benar-benar berfungsi; jangan menyajikan kontrol dekoratif yang tidak memiliki aksi.
- Gunakan URL/rute dan nama aksi yang konsisten dengan proyek. Dokumen ini tidak menetapkan endpoint API.
- Untuk upload gambar, tampilkan progress/error dan preview. Validasi file di sisi server juga diperlukan; validasi UI saja tidak cukup.
- Untuk data transaksi, saldo poin, otorisasi, kode konfirmasi, dan moderasi, backend menjadi sumber kebenaran.

## 14. Data contoh dan aturan placeholder

Data yang muncul pada PDF—misalnya Rusdi, Rani Sari, `@masrusdi`, produk “Kemeja Flanel Kotak”, harga 65K, rating 4,8, jarak 0,8 km, persentase kecocokan 89%, kode contoh, dan jumlah poin—adalah ilustrasi desain.

- Jangan memasukkannya ke kode produksi sebagai nilai default.
- Jika dibutuhkan untuk storybook/demo, letakkan dalam fixture terpisah dan beri label sebagai data demo.
- Jangan gunakan contoh kode transaksi sebagai kode nyata.
- Jangan mengarang foto, skor kondisi, lokasi, rating, atau metrik dampak seolah berasal dari pengguna sungguhan.

## 15. Keputusan yang belum ditetapkan dan perlu dikonfirmasi

Jangan diam-diam mengambil keputusan final untuk hal-hal berikut:

- Apakah “Powered by GROK” tetap menjadi bagian splash/brand final.
- Definisi dan rentang skor kondisi pakaian, serta apakah penilaiannya manual, berbasis aturan, atau dibantu AI.
- Rumus persentase kecocokan ukuran dan ukuran minimum yang diperlukan untuk menghitungnya.
- Aturan level poin selain “Pemula Hijau”, ambang level, kedaluwarsa poin, serta pembatalan/retur.
- Alur sengketa jika kode konfirmasi salah, waktu habis, atau satu pihak tidak merespons.
- Aturan eligibility rating/ulasan, edit/hapus ulasan, dan moderasi ulasan.
- Provider peta, sumber titik donasi, mekanisme lokasi, dan seberapa presisi lokasi yang dapat dilihat publik.
- Payment gateway yang akan dipakai. Desain PDF menunjukkan Cash dan Transfer ke Penjual, sedangkan PRD menyebut payment gateway secara umum; jangan mengklaim integrasi aktif sebelum diputuskan.
- Ketersediaan login Google/Apple, lampiran chat, notifikasi, retur, pengiriman, dan tracking pengiriman.
- Detail peran/permission admin, prosedur moderasi, suspend akun, banding, dan audit.
- Apakah pengguna dapat menyimpan draft listing, mengedit listing aktif, atau menonaktifkan barang terjual.
- Syarat program donasi, rework loop, gratis ongkir, dan tanam pohon; jangan menjanjikan reward yang belum benar-benar tersedia.

## 16. Checklist sebelum UI dianggap selesai

- [ ] Seluruh teks UI menggunakan Bahasa Indonesia yang konsisten dan ejaan yang benar.
- [ ] Desain mengikuti token warna, radius, spasi, dan tipografi.
- [ ] Layar yang tersedia pada PDF dipertahankan; kebutuhan PRD yang belum tergambar sudah ditambahkan dengan pola yang konsisten.
- [ ] Alur pembeli dan penjual tersedia pada akun yang sama.
- [ ] Katalog, pencarian, detail produk, preferensi ukuran, penanda cacat, dan skor kondisi saling terhubung.
- [ ] Checkout menampilkan lokasi, metode pembayaran yang benar-benar tersedia, dan total biaya secara jelas.
- [ ] Konfirmasi transaksi dua arah memiliki state menunggu, gagal, kedaluwarsa, dan sukses sesuai kemampuan backend.
- [ ] Poin hanya tampil sebagai didapat setelah transaksi sah; penukaran memiliki validasi dan konfirmasi.
- [ ] Peta membedakan produk dan titik donasi serta menangani izin lokasi dan error.
- [ ] Chat, ulasan, pelaporan, dan profil memiliki empty/loading/error/success state.
- [ ] Dashboard admin terpisah dan mendukung moderasi serta pengawasan transaksi.
- [ ] Tidak ada data contoh yang bocor menjadi data produksi.
- [ ] Semua kontrol interaktif benar-benar berfungsi dan memiliki feedback.
- [ ] Tampilan mobile, desktop admin, keyboard, focus state, kontras, dan safe area telah diperiksa.
- [ ] Keputusan produk yang belum ditetapkan telah ditandai, bukan ditebak.
