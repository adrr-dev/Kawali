package reviews

import (
	"errors"

	"gorm.io/gorm"

	"mewali/internal/models"
)

var (
	ErrNotDone       = errors.New("hanya transaksi selesai yang bisa diulas")
	ErrNotBuyer      = errors.New("hanya pembeli yang bisa mengulas")
	ErrBadStars      = errors.New("bintang harus 1-5")
	ErrAlreadyExists = errors.New("order ini sudah diulas")
)

// CreateReview membuat ulasan. Syarat: order done + reviewer adalah buyer.
func CreateReview(gdb *gorm.DB, orderID, reviewerID uint, stars int, body string) error {
	if stars < 1 || stars > 5 {
		return ErrBadStars
	}
	var o models.Order
	if err := gdb.First(&o, orderID).Error; err != nil {
		return err
	}
	if o.Status != models.OrderDone {
		return ErrNotDone
	}
	if o.BuyerID != reviewerID {
		return ErrNotBuyer
	}
	var n int64
	gdb.Model(&models.Review{}).Where("order_id = ?", orderID).Count(&n)
	if n > 0 {
		return ErrAlreadyExists
	}
	return gdb.Create(&models.Review{OrderID: orderID, ReviewerID: reviewerID, SellerID: o.SellerID, Stars: stars, Body: body}).Error
}

// SellerRating menghitung rata-rata dan jumlah ulasan penjual.
// Mengembalikan rated=false bila belum ada ulasan (UI: "Belum ada rating").
func SellerRating(gdb *gorm.DB, sellerID uint) (avg float64, count int64, rated bool) {
	type row struct {
		Avg float64
		N   int64
	}
	var r row
	gdb.Model(&models.Review{}).Where("seller_id = ?", sellerID).Select("COALESCE(AVG(stars),0) AS avg, COUNT(*) AS n").Scan(&r)
	return r.Avg, r.N, r.N > 0
}
