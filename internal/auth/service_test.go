package auth

import (
	"path/filepath"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func openSeeded(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "auth.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err := gdb.Create(&models.User{Name: "T", Email: "t@test.id", Handle: "t_1", PasswordHash: string(hash)}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	return gdb
}

func TestLoginOK(t *testing.T) {
	u, err := Login(openSeeded(t), "t@test.id", "password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if u.Email != "t@test.id" {
		t.Errorf("email = %q", u.Email)
	}
}

func TestLoginFails(t *testing.T) {
	gdb := openSeeded(t)
	for name, pair := range map[string][2]string{
		"sandi salah": {"t@test.id", "salah"},
		"tak dikenal": {"x@test.id", "password"},
		"kosong":      {"", ""},
	} {
		if _, err := Login(gdb, pair[0], pair[1]); err != ErrInvalidLogin {
			t.Errorf("%s: mau ErrInvalidLogin, dapat %v", name, err)
		}
	}
}
