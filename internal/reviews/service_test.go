package reviews

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func openTest(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "rev.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	gdb.Create(&models.Order{BuyerID: 2, SellerID: 1, ProductID: 9, Total: 1, PayMethod: models.PayCash, Mode: models.ModeP2P, Status: models.OrderDone})
	gdb.Create(&models.Order{BuyerID: 2, SellerID: 1, ProductID: 8, Total: 1, PayMethod: models.PayCash, Mode: models.ModeP2P, Status: models.OrderPending})
	return gdb
}

func TestCreateReviewOK(t *testing.T) {
	if err := CreateReview(openTest(t), 1, 2, 5, "Bagus"); err != nil {
		t.Fatalf("review: %v", err)
	}
}

func TestReviewGuards(t *testing.T) {
	gdb := openTest(t)
	if err := CreateReview(gdb, 2, 2, 5, "x"); err != ErrNotDone {
		t.Errorf("pending mau ErrNotDone, dapat %v", err)
	}
	if err := CreateReview(gdb, 1, 1, 5, "x"); err != ErrNotBuyer {
		t.Errorf("seller mau ErrNotBuyer, dapat %v", err)
	}
	if err := CreateReview(gdb, 1, 2, 0, "x"); err != ErrBadStars {
		t.Errorf("bintang 0 mau ErrBadStars, dapat %v", err)
	}
	_ = CreateReview(gdb, 1, 2, 4, "ok")
	if err := CreateReview(gdb, 1, 2, 5, "lagi"); err != ErrAlreadyExists {
		t.Errorf("ganda mau ErrAlreadyExists, dapat %v", err)
	}
}

func TestSellerRating(t *testing.T) {
	gdb := openTest(t)
	if _, _, rated := SellerRating(gdb, 1); rated {
		t.Error("tanpa ulasan harus rated=false")
	}
	_ = CreateReview(gdb, 1, 2, 4, "ok")
	avg, n, rated := SellerRating(gdb, 1)
	if !rated || n != 1 || avg != 4 {
		t.Errorf("rating = %v,%d,%v", avg, n, rated)
	}
}
