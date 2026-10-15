package catalog

import (
	"strings"

	"gorm.io/gorm"

	"mewali/internal/models"
)

// Search mencari produk aktif: filter kategori + kata kunci nama/deskripsi.
// LIKE via LOWER agar jalan di SQLite maupun Postgres.
func Search(gdb *gorm.DB, category, q string, limit int) ([]models.Product, error) {
	if limit <= 0 || limit > 60 {
		limit = 24
	}
	tx := gdb.Where("status = ?", models.ProductActive)
	if category != "" {
		tx = tx.Where("category = ?", category)
	}
	if strings.TrimSpace(q) != "" {
		like := "%" + strings.ToLower(strings.TrimSpace(q)) + "%"
		tx = tx.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", like, like)
	}
	var out []models.Product
	err := tx.Order("id DESC").Limit(limit).Find(&out).Error
	return out, err
}

// Cover mengambil foto pertama produk untuk kartu/galeri.
func Cover(gdb *gorm.DB, productID uint) (models.ProductPhoto, error) {
	var ph models.ProductPhoto
	err := gdb.Where("product_id = ?", productID).Order("sort ASC, id ASC").First(&ph).Error
	return ph, err
}

// GetSizes mengambil ukuran barang (bisa kosong bila belum diukur).
func GetSizes(gdb *gorm.DB, productID uint) models.ProductSize {
	var s models.ProductSize
	gdb.Where("product_id = ?", productID).First(&s)
	return s
}
