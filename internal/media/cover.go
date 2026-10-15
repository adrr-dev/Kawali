package media

import (
	"gorm.io/gorm"

	"mewali/internal/models"
)

// Fallback saat produk belum punya foto.
const CoverFallback = "/static/img/placeholder.svg"

// CoverURL mengembalikan foto pertama produk atau fallback.
// Rumah tunggal helper cover agar handler tidak duplikasi/siklus impor.
func CoverURL(gdb *gorm.DB, productID uint) string {
	var ph models.ProductPhoto
	if err := gdb.Where("product_id = ?", productID).Order("sort ASC, id ASC").First(&ph).Error; err == nil && ph.URL != "" {
		return ph.URL
	}
	return CoverFallback
}
