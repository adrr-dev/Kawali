package users

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func TestUpsertPreference(t *testing.T) {
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "pref.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := gdb.Create(&models.User{Name: "U", Email: "u@t.id", Handle: "u_1"}).Error; err != nil {
		t.Fatalf("user: %v", err)
	}
	if err := UpsertPreference(gdb, 1, models.SizePreference{Dada: 52}); err != nil {
		t.Fatalf("insert: %v", err)
	}
	if err := UpsertPreference(gdb, 1, models.SizePreference{Dada: 54, Torso: 71}); err != nil {
		t.Fatalf("update: %v", err)
	}
	got := GetPreference(gdb, 1)
	if got.Dada != 54 || got.Torso != 71 {
		t.Errorf("pref = %+v", got)
	}
	if err := UpdateLocation(gdb, 1, "Jakarta", -6.2, 106.8); err != nil {
		t.Fatalf("loc: %v", err)
	}
	var u models.User
	gdb.First(&u, 1)
	if u.Lat != -6.2 || u.LocationLabel != "Jakarta" {
		t.Errorf("user = %+v", u)
	}
}
