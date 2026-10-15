package orders

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"mewali/internal/models"
	"mewali/internal/points"
)

var (
	ErrNotActive    = errors.New("produk tidak tersedia")
	ErrSelfOrder    = errors.New("tidak bisa membeli barang sendiri")
	ErrBadPayMethod = errors.New("metode bayar harus cash atau transfer_manual")
	ErrBadMode      = errors.New("mode harus p2p atau meetup")
	ErrOrderBusy    = errors.New("produk sedang dalam transaksi")
	ErrCodeWrong    = errors.New("kode salah")
	ErrCodeExpired  = errors.New("kode kedaluwarsa")
	ErrCodeUsed     = errors.New("kode sudah dipakai")
	ErrWrongMode    = errors.New("aksi tidak sesuai mode order")
)

// AddToCart memasukkan produk aktif ke keranjang (satu per user+produk).
func AddToCart(gdb *gorm.DB, userID, productID uint) error {
	var p models.Product
	if err := gdb.First(&p, productID).Error; err != nil {
		return err
	}
	if p.Status != models.ProductActive {
		return ErrNotActive
	}
	return gdb.FirstOrCreate(&models.CartItem{}, models.CartItem{UserID: userID, ProductID: productID}).Error
}

// RemoveFromCart mengeluarkan item dari keranjang.
func RemoveFromCart(gdb *gorm.DB, userID, productID uint) error {
	return gdb.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.CartItem{}).Error
}

// Checkout membuat order pending. Mode meetup langsung menerbitkan 2 kode;
// mode p2p tanpa kode (selesai manual). Mengembalikan kode milik pembeli dan
// penjual (kosong untuk p2p).
func Checkout(gdb *gorm.DB, buyerID, productID uint, payMethod, mode string, now time.Time) (models.Order, string, string, error) {
	var nilOrder models.Order
	if payMethod != models.PayCash && payMethod != models.PayTransferManual {
		return nilOrder, "", "", ErrBadPayMethod
	}
	if mode != models.ModeP2P && mode != models.ModeMeetup {
		return nilOrder, "", "", ErrBadMode
	}
	var p models.Product
	if err := gdb.First(&p, productID).Error; err != nil {
		return nilOrder, "", "", err
	}
	if p.Status != models.ProductActive {
		return nilOrder, "", "", ErrNotActive
	}
	if p.SellerID == buyerID {
		return nilOrder, "", "", ErrSelfOrder
	}
	var busy int64
	gdb.Model(&models.Order{}).Where("product_id = ? AND status = ?", productID, models.OrderPending).Count(&busy)
	if busy > 0 {
		return nilOrder, "", "", ErrOrderBusy
	}

	var order models.Order
	err := gdb.Transaction(func(tx *gorm.DB) error {
		order = models.Order{BuyerID: buyerID, SellerID: p.SellerID, ProductID: productID,
			Total: p.Price, PayMethod: payMethod, Mode: mode, Status: models.OrderPending}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		tx.Where("user_id = ? AND product_id = ?", buyerID, productID).Delete(&models.CartItem{})
		return nil
	})
	if err != nil {
		return nilOrder, "", "", err
	}

	if mode == models.ModeP2P {
		return order, "", "", nil
	}
	buyerCode, err := Generate()
	if err != nil {
		return nilOrder, "", "", err
	}
	sellerCode, err := Generate()
	if err != nil {
		return nilOrder, "", "", err
	}
	exp := ExpiredAt(now)
	rows := []models.ConfirmCode{
		{OrderID: order.ID, Role: models.CodeBuyer, Code: buyerCode, ExpiresAt: exp},
		{OrderID: order.ID, Role: models.CodeSeller, Code: sellerCode, ExpiresAt: exp},
	}
	if err := gdb.Create(&rows).Error; err != nil {
		return nilOrder, "", "", err
	}
	return order, buyerCode, sellerCode, nil
}

// counterpart menukar peran.
func counterpart(role string) string {
	if role == models.CodeBuyer {
		return models.CodeSeller
	}
	return models.CodeBuyer
}

// Confirm memasukkan kode MILIK LAWAN. Saat kedua kode terpakai, order done,
// produk sold, pembeli dapat poin (mode meetup saja).
func Confirm(gdb *gorm.DB, orderID uint, submitterRole, code string, now time.Time) error {
	var order models.Order
	if err := gdb.First(&order, orderID).Error; err != nil {
		return err
	}
	if order.Mode != models.ModeMeetup {
		return ErrWrongMode
	}
	if order.Status != models.OrderPending {
		return ErrCodeUsed
	}
	var row models.ConfirmCode
	if err := gdb.Where("order_id = ? AND role = ?", orderID, counterpart(submitterRole)).First(&row).Error; err != nil {
		return ErrCodeWrong
	}
	if row.UsedAt != nil {
		return ErrCodeUsed
	}
	if IsExpired(row.ExpiresAt, now) {
		return ErrCodeExpired
	}
	if row.Code != code {
		gdb.Model(&row).Update("attempts", row.Attempts+1)
		return ErrCodeWrong
	}
	return gdb.Transaction(func(tx *gorm.DB) error {
		t := now
		if err := tx.Model(&row).Update("used_at", &t).Error; err != nil {
			return err
		}
		var used int64
		if err := tx.Model(&models.ConfirmCode{}).Where("order_id = ? AND used_at IS NOT NULL", orderID).Count(&used).Error; err != nil {
			return err
		}
		if used < 2 {
			return nil
		}
		if err := tx.Model(&order).Update("status", models.OrderDone).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Product{}).Where("id = ?", order.ProductID).Update("status", models.ProductSold).Error; err != nil {
			return err
		}
		delta, _ := points.EarnFor(points.ActionBeli)
		var balance int64
		tx.Model(&models.PointLedger{}).Where("user_id = ?", order.BuyerID).Select("COALESCE(SUM(delta),0)").Scan(&balance)
		oid := order.ID
		return tx.Create(&models.PointLedger{UserID: order.BuyerID, OrderID: &oid,
			Delta: delta, Reason: string(points.ActionBeli), BalanceAfter: int(balance) + delta}).Error
	})
}

// P2PDone menyelesaikan order mode p2p secara manual (tanpa kode/poin otomatis).
func P2PDone(gdb *gorm.DB, orderID uint) error {
	var order models.Order
	if err := gdb.First(&order, orderID).Error; err != nil {
		return err
	}
	if order.Mode != models.ModeP2P {
		return ErrWrongMode
	}
	if order.Status != models.OrderPending {
		return ErrCodeUsed
	}
	return gdb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&order).Update("status", models.OrderDone).Error; err != nil {
			return err
		}
		return tx.Model(&models.Product{}).Where("id = ?", order.ProductID).Update("status", models.ProductSold).Error
	})
}
