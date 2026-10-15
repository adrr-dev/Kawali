package models

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "models.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := gdb.AutoMigrate(All()...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return gdb
}

func TestMigrateAllTables(t *testing.T) {
	gdb := openTestDB(t)
	for _, m := range All() {
		if !gdb.Migrator().HasTable(m) {
			t.Errorf("tabel hilang untuk %T", m)
		}
	}
}

func TestUserEmailUnique(t *testing.T) {
	gdb := openTestDB(t)
	u := User{Name: "A", Email: "a@test.id", Handle: "a_1"}
	if err := gdb.Create(&u).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	dup := User{Name: "B", Email: "a@test.id", Handle: "b_1"}
	if err := gdb.Create(&dup).Error; err == nil {
		t.Error("email duplikat harus ditolak unique index")
	}
}

func TestCartItemOnePerUserProduct(t *testing.T) {
	gdb := openTestDB(t)
	if err := gdb.Create(&CartItem{UserID: 1, ProductID: 2}).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := gdb.Create(&CartItem{UserID: 1, ProductID: 2}).Error; err == nil {
		t.Error("cart ganda user+produk harus ditolak")
	}
}

func TestOrderAndCodeFlow(t *testing.T) {
	gdb := openTestDB(t)
	o := Order{BuyerID: 1, SellerID: 2, ProductID: 3, Total: 65000,
		PayMethod: PayCash, Mode: ModeMeetup, Status: OrderPending}
	if err := gdb.Create(&o).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}
	for _, role := range []string{CodeBuyer, CodeSeller} {
		c := ConfirmCode{OrderID: o.ID, Role: role, Code: "ABCDE"}
		if err := gdb.Create(&c).Error; err != nil {
			t.Fatalf("create code %s: %v", role, err)
		}
	}
	dup := ConfirmCode{OrderID: o.ID, Role: CodeBuyer, Code: "XYZ12"}
	if err := gdb.Create(&dup).Error; err == nil {
		t.Error("kode ganda order+role harus ditolak")
	}
}

func TestReviewOnePerOrder(t *testing.T) {
	gdb := openTestDB(t)
	r := Review{OrderID: 9, ReviewerID: 1, SellerID: 2, Stars: 5, Body: "Bagus"}
	if err := gdb.Create(&r).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	dup := Review{OrderID: 9, ReviewerID: 1, SellerID: 2, Stars: 1, Body: "Jelek"}
	if err := gdb.Create(&dup).Error; err == nil {
		t.Error("ulasan kedua untuk order sama harus ditolak")
	}
}
