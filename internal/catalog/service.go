package catalog

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mewali/internal/models"
)

// Slot foto wajib sesuai DESIGN Bab tambah-katalog.
var RequiredSlots = []string{"depan", "belakang", "kerah", "lengan", "tag_samping", "tag_atas"}

func validSlot(s string) bool {
	for _, r := range RequiredSlots {
		if r == s {
			return true
		}
	}
	return false
}

// CreateProduct membuat listing status draft.
func CreateProduct(gdb *gorm.DB, sellerID uint, name, category string, price int) (models.Product, error) {
	if name == "" || price <= 0 {
		return models.Product{}, errors.New("nama dan harga (>0) wajib diisi")
	}
	p := models.Product{SellerID: sellerID, Name: name, Category: category, Price: price, Status: models.ProductDraft}
	return p, gdb.Create(&p).Error
}

// AddPhoto menambah foto ke slot. Slot di luar daftar wajib ditolak.
func AddPhoto(gdb *gorm.DB, productID uint, slot, url string) error {
	if !validSlot(slot) {
		return fmt.Errorf("slot %q tidak dikenal", slot)
	}
	if url == "" {
		return errors.New("url foto wajib diisi")
	}
	return gdb.Create(&models.ProductPhoto{ProductID: productID, Slot: slot, URL: url}).Error
}

// FilledSlots mengembalikan slot yang sudah terisi (unik) untuk produk.
func FilledSlots(gdb *gorm.DB, productID uint) (map[string]bool, error) {
	var slots []string
	if err := gdb.Model(&models.ProductPhoto{}).Where("product_id = ?", productID).Distinct().Pluck("slot", &slots).Error; err != nil {
		return nil, err
	}
	m := map[string]bool{}
	for _, s := range slots {
		m[s] = true
	}
	return m, nil
}

// MissingSlots mendaftar slot wajib yang belum ada fotonya (feedback per slot).
func MissingSlots(filled map[string]bool) []string {
	var miss []string
	for _, r := range RequiredSlots {
		if !filled[r] {
			miss = append(miss, r)
		}
	}
	return miss
}

// AddDefect menandai cacat pada foto (koordinat 0-1000, lihat ARCHITECTURE Bab 3).
func AddDefect(gdb *gorm.DB, productID, photoID uint, x, y int, label string) error {
	if label == "" {
		return errors.New("deskripsi cacat wajib diisi")
	}
	if x < 0 || x > 1000 || y < 0 || y > 1000 {
		return errors.New("koordinat marker harus 0-1000")
	}
	return gdb.Create(&models.Defect{ProductID: productID, PhotoID: photoID, X: x, Y: y, Label: label}).Error
}

// Publish menerbitkan draft: wajib 6 slot + hitung skor kondisi. Gagal
// mengembalikan daftar slot yang kurang agar UI bisa feedback per slot.
func Publish(gdb *gorm.DB, productID uint) (models.Product, error) {
	var p models.Product
	if err := gdb.First(&p, productID).Error; err != nil {
		return p, err
	}
	filled, err := FilledSlots(gdb, productID)
	if err != nil {
		return p, err
	}
	if miss := MissingSlots(filled); len(miss) > 0 {
		return p, fmt.Errorf("foto kurang: %v", miss)
	}
	var defects int64
	if err := gdb.Model(&models.Defect{}).Where("product_id = ?", productID).Count(&defects).Error; err != nil {
		return p, err
	}
	score, _ := Score(len(RequiredSlots), int(defects))
	now := time.Now()
	cs := models.ConditionScore{ProductID: productID, Score: score, Method: MethodRule, CheckedAt: now}
	if err := gdb.Where("product_id = ?", productID).FirstOrCreate(&cs).Error; err != nil {
		return p, err
	}
	gdb.Model(&cs).Updates(map[string]any{"score": score, "method": MethodRule, "checked_at": now})
	p.Status = models.ProductActive
	return p, gdb.Save(&p).Error
}

// SaveSizes menyimpan (upsert) ukuran barang dalam cm.
func SaveSizes(gdb *gorm.DB, productID uint, s models.ProductSize) error {
	s.ProductID = productID
	var cur models.ProductSize
	if err := gdb.Where("product_id = ?", productID).First(&cur).Error; err != nil {
		return gdb.Create(&s).Error
	}
	s.ID = cur.ID
	return gdb.Save(&s).Error
}

// SellerList mendaftar produk milik penjual (termasuk draft) untuk kelola.
func SellerList(gdb *gorm.DB, sellerID uint) ([]models.Product, error) {
	var out []models.Product
	err := gdb.Where("seller_id = ?", sellerID).Order("id DESC").Find(&out).Error
	return out, err
}

// PhotosBySlot mengelompokkan foto per slot untuk checklist UI.
func PhotosBySlot(gdb *gorm.DB, productID uint) (map[string][]models.ProductPhoto, error) {
	var all []models.ProductPhoto
	if err := gdb.Where("product_id = ?", productID).Order("id ASC").Find(&all).Error; err != nil {
		return nil, err
	}
	m := map[string][]models.ProductPhoto{}
	for _, ph := range all {
		m[ph.Slot] = append(m[ph.Slot], ph)
	}
	return m, nil
}
