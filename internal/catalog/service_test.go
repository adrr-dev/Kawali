package catalog

import (
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func openTest(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "cat.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return gdb
}

func fillAllSlots(t *testing.T, gdb *gorm.DB, id uint) {
	t.Helper()
	for i, s := range RequiredSlots {
		if err := AddPhoto(gdb, id, s, "http://x/"+s+".jpg"); err != nil {
			t.Fatalf("slot %s: %v", s, err)
		}
		_ = i
	}
}

func TestPublishNeedsAllSlots(t *testing.T) {
	gdb := openTest(t)
	p, err := CreateProduct(gdb, 1, "Kemeja", "Kemeja", 65000)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := Publish(gdb, p.ID); err == nil {
		t.Fatal("tanpa foto harus gagal")
	}
	if err := AddPhoto(gdb, p.ID, "depan", "http://x/d.jpg"); err != nil {
		t.Fatalf("add: %v", err)
	}
	err = func() error {
		_, err := Publish(gdb, p.ID)
		return err
	}()
	if err == nil || !strings.Contains(err.Error(), "foto kurang") {
		t.Fatalf("mau error foto kurang, dapat %v", err)
	}
}

func TestPublishOKScoresWithDefects(t *testing.T) {
	gdb := openTest(t)
	p, _ := CreateProduct(gdb, 1, "Jeans", "Jeans", 40000)
	fillAllSlots(t, gdb, p.ID)
	if err := AddDefect(gdb, p.ID, 1, 100, 200, "Pudar di kerah"); err != nil {
		t.Fatalf("defect: %v", err)
	}
	got, err := Publish(gdb, p.ID)
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if got.Status != models.ProductActive {
		t.Errorf("status = %q", got.Status)
	}
	var cs models.ConditionScore
	gdb.Where("product_id = ?", p.ID).First(&cs)
	if cs.Score != 90 || cs.Method != MethodRule {
		t.Errorf("skor = %d/%q, mau 90/aturan", cs.Score, cs.Method)
	}
}

func TestInvalidInputsRejected(t *testing.T) {
	gdb := openTest(t)
	if _, err := CreateProduct(gdb, 1, "", "Kaos", 0); err == nil {
		t.Error("produk kosong harus ditolak")
	}
	p, _ := CreateProduct(gdb, 1, "Kaos", "Kaos", 30000)
	if err := AddPhoto(gdb, p.ID, "samping", "http://x.jpg"); err == nil {
		t.Error("slot asing harus ditolak")
	}
	if err := AddDefect(gdb, p.ID, 1, 2000, 0, "x"); err == nil {
		t.Error("koordinat di luar 0-1000 harus ditolak")
	}
	if err := AddDefect(gdb, p.ID, 1, 10, 10, ""); err == nil {
		t.Error("label kosong harus ditolak")
	}
}

func TestSaveSizesUpsert(t *testing.T) {
	gdb := openTest(t)
	p, _ := CreateProduct(gdb, 1, "X", "Kaos", 10000)
	if err := SaveSizes(gdb, p.ID, struct2size(50, 70)); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := SaveSizes(gdb, p.ID, struct2size(54, 71)); err != nil {
		t.Fatalf("resave: %v", err)
	}
	var n int64
	gdb.Model(&models.ProductSize{}).Where("product_id = ?", p.ID).Count(&n)
	if n != 1 {
		t.Errorf("rows = %d, mau 1 (upsert)", n)
	}
	got := GetSizes(gdb, p.ID)
	if got.Dada != 54 || got.Torso != 71 {
		t.Errorf("sizes = %+v", got)
	}
}

func struct2size(dada, torso float64) models.ProductSize {
	return models.ProductSize{Dada: dada, Torso: torso}
}

func TestSellerList(t *testing.T) {
	gdb := openTest(t)
	p, _ := CreateProduct(gdb, 7, "A", "Kaos", 10000)
	items, err := SellerList(gdb, 7)
	if err != nil || len(items) != 1 || items[0].ID != p.ID {
		t.Errorf("list = %+v, err = %v", items, err)
	}
}
