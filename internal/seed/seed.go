package seed

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mewali/internal/catalog"
	"mewali/internal/models"
)

// Paket sendiri (bukan db) agar seed demo boleh memakai service catalog
// tanpa siklus impor dengan test yang memakai db.Migrate.

type filler struct {
	name   string
	email  string
	handle string
	level  string
	admin  bool
}

// Run membuat akun filler + satu listing demo lewat jalur Publish asli.
// Idempoten: aman dijalankan ulang.
func Run(gdb *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		return err
	}
	for _, f := range []filler{
		{"Pembeli Demo", "pembeli@test.id", "pembeli_demo", "Pemula Hijau", false},
		{"Penjual Demo", "penjual@test.id", "penjual_demo", "Pemula Hijau", false},
		{"Admin Demo", "admin@test.id", "admin_demo", "Pemula Hijau", true},
	} {
		u := models.User{Name: f.name, Email: f.email, PasswordHash: string(hash), Handle: f.handle, Level: f.level, IsAdmin: f.admin}
		if err := gdb.Where("email = ?", f.email).FirstOrCreate(&u).Error; err != nil {
			return err
		}
		// pastikan flag admin benar walau baris sudah ada sebelum kolomnya ada
		if f.admin && !u.IsAdmin {
			gdb.Model(&u).Update("is_admin", true)
		}
	}
	if err := donationPoints(gdb); err != nil {
		return err
	}
	if err := demoBanners(gdb); err != nil {
		return err
	}
	return demoProduct(gdb)
}

func donationPoints(gdb *gorm.DB) error {
	for _, d := range []models.DonationPoint{
		{Name: "Titik Donasi Demo 1", Lat: -6.21, Lng: 106.81, AddressLabel: "Jakarta Selatan", Active: true},
		{Name: "Titik Donasi Demo 2", Lat: -6.19, Lng: 106.83, AddressLabel: "Jakarta Pusat", Active: true},
	} {
		var n int64
		gdb.Model(&models.DonationPoint{}).Where("name = ?", d.Name).Count(&n)
		if n == 0 {
			if err := gdb.Create(&d).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func demoBanners(gdb *gorm.DB) error {
	for i, b := range []models.Banner{
		{Title: "Terdekat dengan dirimu", Subtitle: "Temukan pakaian di sekitarmu.", CTALabel: "Beli", CTAHref: "/cari", Sort: 1, Active: true},
		{Title: "Raih Poin", Subtitle: "Tukarkan poinmu menjadi sesuatu yang bermakna.", CTALabel: "Lihat Poin", CTAHref: "/poin", Sort: 2, Active: true},
		{Title: "Jangan Akhiri Ceritanya di Lemari", Subtitle: "Putar kembali siklus pakaianmu.", CTALabel: "Mulai Jual", CTAHref: "/tambah", Sort: 3, Active: true},
	} {
		var n int64
		gdb.Model(&models.Banner{}).Where("title = ?", b.Title).Count(&n)
		if n == 0 {
			b.Sort = i + 1
			if err := gdb.Create(&b).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func demoProduct(gdb *gorm.DB) error {
	var seller models.User
	if err := gdb.Where("email = ?", "penjual@test.id").First(&seller).Error; err != nil {
		return err
	}
	var n int64
	gdb.Model(&models.Product{}).Where("seller_id = ? AND name = ?", seller.ID, "Kemeja Flanel Demo").Count(&n)
	if n > 0 {
		return nil
	}
	p, err := catalog.CreateProduct(gdb, seller.ID, "Kemeja Flanel Demo", "Kemeja", 65000)
	if err != nil {
		return err
	}
	for _, s := range catalog.RequiredSlots {
		if err := catalog.AddPhoto(gdb, p.ID, s, "/static/img/placeholder.svg"); err != nil {
			return err
		}
	}
	size := models.ProductSize{ProductID: p.ID, Dada: 54, Pinggang: 80, Pinggul: 90, Tinggi: 170, Torso: 71}
	if err := gdb.Create(&size).Error; err != nil {
		return err
	}
	_, err = catalog.Publish(gdb, p.ID)
	return err
}
