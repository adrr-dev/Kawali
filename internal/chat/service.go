package chat

import (
	"errors"

	"gorm.io/gorm"

	"mewali/internal/models"
)

var (
	ErrNotMember = errors.New("bukan anggota thread ini")
	ErrEmptyBody = errors.New("pesan kosong")
)

// GetOrCreateThread membuka atau memakai ulang thread produk antara pembeli-penjual.
func GetOrCreateThread(gdb *gorm.DB, productID, buyerID, sellerID uint, orderID *uint) (models.ChatThread, error) {
	var th models.ChatThread
	q := gdb.Where("product_id = ? AND buyer_id = ? AND seller_id = ?", productID, buyerID, sellerID)
	if orderID != nil {
		q = q.Where("order_id = ?", *orderID)
	}
	if err := q.First(&th).Error; err == nil {
		return th, nil
	}
	th = models.ChatThread{OrderID: orderID, ProductID: productID, BuyerID: buyerID, SellerID: sellerID}
	return th, gdb.Create(&th).Error
}

// PostMessage mengirim pesan. Pengirim harus buyer/seller thread.
func PostMessage(gdb *gorm.DB, threadID, senderID uint, body string) error {
	if body == "" {
		return ErrEmptyBody
	}
	var th models.ChatThread
	if err := gdb.First(&th, threadID).Error; err != nil {
		return err
	}
	if senderID != th.BuyerID && senderID != th.SellerID {
		return ErrNotMember
	}
	return gdb.Create(&models.ChatMessage{ThreadID: threadID, SenderID: senderID, Body: body}).Error
}

// ListMessages mengembalikan pesan terlama-dulu, dibatasi limit.
func ListMessages(gdb *gorm.DB, threadID uint, limit int) ([]models.ChatMessage, error) {
	var out []models.ChatMessage
	err := gdb.Where("thread_id = ?", threadID).Order("id ASC").Limit(limit).Find(&out).Error
	return out, err
}
