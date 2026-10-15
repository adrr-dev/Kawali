package catalog

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func openSearch(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "s.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	gdb.Create(&models.Product{SellerID: 1, Name: "Kemeja Flanel", Category: "Kemeja", Price: 65000, Status: models.ProductActive})
	gdb.Create(&models.Product{SellerID: 1, Name: "Jeans Sobek", Category: "Jeans", Price: 40000, Status: models.ProductActive})
	gdb.Create(&models.Product{SellerID: 1, Name: "Draft Rahasia", Category: "Kaos", Price: 10000, Status: models.ProductDraft})
	return gdb
}

func TestSearchFilters(t *testing.T) {
	gdb := openSearch(t)
	all, _ := Search(gdb, "", "", 24)
	if len(all) != 2 {
		t.Fatalf("semua = %d, mau 2 (draft dikecualikan)", len(all))
	}
	kat, _ := Search(gdb, "Jeans", "", 24)
	if len(kat) != 1 || kat[0].Name != "Jeans Sobek" {
		t.Errorf("kategori salah: %+v", kat)
	}
	q, _ := Search(gdb, "", "flanel", 24)
	if len(q) != 1 {
		t.Errorf("keyword = %d, mau 1", len(q))
	}
	q2, _ := Search(gdb, "", "KEMEJA", 24)
	if len(q2) != 1 {
		t.Errorf("case-insensitive = %d, mau 1", len(q2))
	}
	none, _ := Search(gdb, "", "gaun pesta", 24)
	if len(none) != 0 {
		t.Errorf("kosong = %d", len(none))
	}
}
