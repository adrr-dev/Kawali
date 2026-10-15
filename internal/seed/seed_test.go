package seed

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func TestRunIdempotent(t *testing.T) {
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "seed.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	for i := 0; i < 2; i++ {
		if err := Run(gdb); err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
	}
	var users int64
	gdb.Model(&models.User{}).Count(&users)
	if users != 3 {
		t.Errorf("users = %d, mau 3", users)
	}
	var p models.Product
	if err := gdb.Where("name = ?", "Kemeja Flanel Demo").First(&p).Error; err != nil {
		t.Fatalf("demo hilang: %v", err)
	}
	if p.Status != models.ProductActive {
		t.Errorf("demo = %q, mau active", p.Status)
	}
}
