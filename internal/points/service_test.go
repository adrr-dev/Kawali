package points

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
)

func openTest(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "pts.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return gdb
}

func TestAwardAndBalance(t *testing.T) {
	gdb := openTest(t)
	if Balance(gdb, 1) != 0 {
		t.Fatal("saldo awal harus 0")
	}
	if err := Award(gdb, 1, nil, ActionBeli); err != nil {
		t.Fatalf("award: %v", err)
	}
	if err := Award(gdb, 1, nil, ActionJualLengkap); err != nil {
		t.Fatalf("award: %v", err)
	}
	if b := Balance(gdb, 1); b != 25 {
		t.Errorf("saldo = %d, mau 25", b)
	}
	if err := Award(gdb, 1, nil, "ngarang"); err == nil {
		t.Error("aksi aneh harus gagal")
	}
}

func TestRedeemServiceFlow(t *testing.T) {
	gdb := openTest(t)
	if err := Redeem(gdb, 1, RewardDonasi); err != ErrInsufficient {
		t.Fatalf("miskin mau ErrInsufficient, dapat %v", err)
	}
	_ = Award(gdb, 1, nil, ActionReferral) // +25
	_ = Award(gdb, 1, nil, ActionReferral)
	_ = Award(gdb, 1, nil, ActionReferral)
	_ = Award(gdb, 1, nil, ActionBeli) // total 85
	if err := Redeem(gdb, 1, RewardDonasi); err != nil {
		t.Fatalf("redeem 80 dari 85: %v", err)
	}
	if b := Balance(gdb, 1); b != 5 {
		t.Errorf("saldo = %d, mau 5", b)
	}
	if err := Redeem(gdb, 1, "ngarang"); err == nil {
		t.Error("reward aneh harus gagal")
	}
}
