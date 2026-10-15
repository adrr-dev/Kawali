package orders

import (
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func openTest(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "ord.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	gdb.Create(&models.User{Name: "S", Email: "s@t.id", Handle: "s_1"})
	gdb.Create(&models.User{Name: "B", Email: "b@t.id", Handle: "b_1"})
	gdb.Create(&models.Product{SellerID: 1, Name: "Kemeja", Category: "Kemeja", Price: 65000, Status: models.ProductActive})
	return gdb
}

func TestCheckoutMeetupIssuesCodes(t *testing.T) {
	gdb := openTest(t)
	o, bc, sc, err := Checkout(gdb, 2, 1, models.PayCash, models.ModeMeetup, time.Now())
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}
	if o.Status != models.OrderPending || o.Total != 65000 {
		t.Errorf("order aneh: %+v", o)
	}
	if !ValidFormat(bc) || !ValidFormat(sc) || bc == sc {
		t.Errorf("kode jelek: %q %q", bc, sc)
	}
}

func TestCheckoutGuards(t *testing.T) {
	gdb := openTest(t)
	now := time.Now()
	if _, _, _, err := Checkout(gdb, 1, 1, models.PayCash, models.ModeMeetup, now); err != ErrSelfOrder {
		t.Errorf("beli sendiri mau ErrSelfOrder, dapat %v", err)
	}
	if _, _, _, err := Checkout(gdb, 2, 1, "qr", models.ModeMeetup, now); err != ErrBadPayMethod {
		t.Errorf("pay aneh mau ErrBadPayMethod, dapat %v", err)
	}
	if _, _, _, err := Checkout(gdb, 2, 1, models.PayCash, models.ModeMeetup, now); err != nil {
		t.Fatalf("checkout1: %v", err)
	}
	if _, _, _, err := Checkout(gdb, 2, 1, models.PayCash, models.ModeMeetup, now); err != ErrOrderBusy {
		t.Errorf("double checkout mau ErrOrderBusy, dapat %v", err)
	}
}

func TestTwoWayConfirmCompletesWithPoints(t *testing.T) {
	gdb := openTest(t)
	now := time.Now()
	o, bc, sc, err := Checkout(gdb, 2, 1, models.PayTransferManual, models.ModeMeetup, now)
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}
	// kode sendiri ditolak (harus kode lawan)
	if err := Confirm(gdb, o.ID, models.CodeBuyer, bc, now); err != ErrCodeWrong {
		t.Fatalf("kode sendiri mau ErrCodeWrong, dapat %v", err)
	}
	if err := Confirm(gdb, o.ID, models.CodeBuyer, sc, now); err != nil {
		t.Fatalf("buyer confirm: %v", err)
	}
	var mid models.Order
	gdb.First(&mid, o.ID)
	if mid.Status != models.OrderPending {
		t.Fatalf("satu sisi belum boleh done: %q", mid.Status)
	}
	if err := Confirm(gdb, o.ID, models.CodeSeller, bc, now); err != nil {
		t.Fatalf("seller confirm: %v", err)
	}
	var done models.Order
	gdb.First(&done, o.ID)
	if done.Status != models.OrderDone {
		t.Errorf("status = %q, mau done", done.Status)
	}
	var p models.Product
	gdb.First(&p, 1)
	if p.Status != models.ProductSold {
		t.Errorf("produk = %q, mau sold", p.Status)
	}
	var led models.PointLedger
	if err := gdb.Where("user_id = ?", 2).First(&led).Error; err != nil {
		t.Fatalf("ledger: %v", err)
	}
	if led.Delta != 10 || led.BalanceAfter != 10 {
		t.Errorf("ledger aneh: %+v", led)
	}
	// konfirmasi ulang ditolak
	if err := Confirm(gdb, o.ID, models.CodeBuyer, sc, now); err == nil {
		t.Error("confirm setelah done harus gagal")
	}
}

func TestConfirmExpired(t *testing.T) {
	gdb := openTest(t)
	past := time.Now().Add(-2 * time.Hour)
	o, _, sc, err := Checkout(gdb, 2, 1, models.PayCash, models.ModeMeetup, past)
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}
	if err := Confirm(gdb, o.ID, models.CodeBuyer, sc, time.Now()); err != ErrCodeExpired {
		t.Errorf("mau ErrCodeExpired, dapat %v", err)
	}
}

func TestP2PFlow(t *testing.T) {
	gdb := openTest(t)
	o, bc, sc, err := Checkout(gdb, 2, 1, models.PayCash, models.ModeP2P, time.Now())
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}
	if bc != "" || sc != "" {
		t.Error("p2p tidak menerbitkan kode")
	}
	if err := Confirm(gdb, o.ID, models.CodeBuyer, "ABCDE", time.Now()); err != ErrWrongMode {
		t.Errorf("confirm di p2p mau ErrWrongMode, dapat %v", err)
	}
	if err := P2PDone(gdb, o.ID); err != nil {
		t.Fatalf("p2p done: %v", err)
	}
	var done models.Order
	gdb.First(&done, o.ID)
	if done.Status != models.OrderDone {
		t.Errorf("status = %q", done.Status)
	}
	var n int64
	gdb.Model(&models.PointLedger{}).Count(&n)
	if n != 0 {
		t.Error("p2p tidak memberi poin otomatis")
	}
}

func TestCartRoundtrip(t *testing.T) {
	gdb := openTest(t)
	if err := AddToCart(gdb, 2, 1); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := AddToCart(gdb, 2, 1); err != nil {
		t.Fatalf("add ganda harus idempoten: %v", err)
	}
	var n int64
	gdb.Model(&models.CartItem{}).Where("user_id = 2").Count(&n)
	if n != 1 {
		t.Errorf("cart = %d, mau 1", n)
	}
	if _, _, _, err := Checkout(gdb, 2, 1, models.PayCash, models.ModeP2P, time.Now()); err != nil {
		t.Fatalf("checkout: %v", err)
	}
	gdb.Model(&models.CartItem{}).Where("user_id = 2").Count(&n)
	if n != 0 {
		t.Error("checkout harus kosongkan item dari cart")
	}
}
