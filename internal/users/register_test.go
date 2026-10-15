package users

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
)

func openReg(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "reg.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return gdb
}

func TestRegisterOK(t *testing.T) {
	u, err := Register(openReg(t), "Rani Sari", "rani@test.id", "rahasia123")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if u.Handle != "ranisari" {
		t.Errorf("handle = %q", u.Handle)
	}
}

func TestRegisterGuards(t *testing.T) {
	gdb := openReg(t)
	if _, err := Register(gdb, "", "a@t.id", "rahasia123"); err != ErrBadSignup {
		t.Errorf("nama kosong: %v", err)
	}
	if _, err := Register(gdb, "A", "a@t.id", "pendek"); err != ErrBadSignup {
		t.Errorf("sandi pendek: %v", err)
	}
	if _, err := Register(gdb, "A", "a@t.id", "rahasia123"); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if _, err := Register(gdb, "B", "a@t.id", "rahasia123"); err != ErrEmailTaken {
		t.Errorf("duplikat: %v", err)
	}
	u, err := Register(gdb, "A", "b@t.id", "rahasia123")
	if err != nil {
		t.Fatalf("nama sama: %v", err)
	}
	if u.Handle != "a2" {
		t.Errorf("handle duplikat = %q, mau a2", u.Handle)
	}
}
