package reports

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
)

func TestFileReport(t *testing.T) {
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "rep.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	uid, pid := uint(5), uint(9)
	if err := FileReport(gdb, 1, &uid, nil, "spam", "detail"); err != nil {
		t.Errorf("lapor user: %v", err)
	}
	if err := FileReport(gdb, 1, nil, &pid, "palsu", ""); err != nil {
		t.Errorf("lapor produk: %v", err)
	}
	if err := FileReport(gdb, 1, nil, nil, "x", ""); err != ErrNoTarget {
		t.Errorf("tanpa target mau ErrNoTarget, dapat %v", err)
	}
	if err := FileReport(gdb, 1, &uid, nil, "", ""); err != ErrNoReason {
		t.Errorf("tanpa alasan mau ErrNoReason, dapat %v", err)
	}
}
