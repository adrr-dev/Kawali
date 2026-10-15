package admin

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
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "adm.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	gdb.Create(&models.Report{ReporterID: 1, Reason: "spam", Status: models.ReportOpen})
	gdb.Create(&models.Product{SellerID: 1, Name: "X", Price: 1, Status: models.ProductActive})
	return gdb
}

func TestOverviewCounts(t *testing.T) {
	o := GetOverview(openTest(t))
	if o.OpenReports != 1 || o.Products != 1 {
		t.Errorf("overview = %+v", o)
	}
}

func TestResolveAndTakedownAudit(t *testing.T) {
	gdb := openTest(t)
	if err := ResolveReport(gdb, 9, 1, true); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	var r models.Report
	gdb.First(&r, 1)
	if r.Status != models.ReportResolved {
		t.Errorf("status = %q", r.Status)
	}
	if err := Takedown(gdb, 9, 1); err != nil {
		t.Fatalf("takedown: %v", err)
	}
	var p models.Product
	gdb.First(&p, 1)
	if p.Status != models.ProductTakenDown {
		t.Errorf("produk = %q", p.Status)
	}
	var n int64
	gdb.Model(&models.AdminAudit{}).Where("admin_id = 9").Count(&n)
	if n != 2 {
		t.Errorf("audit = %d, mau 2", n)
	}
	if err := ResolveReport(gdb, 9, 999, true); err == nil {
		t.Error("laporan gaib harus gagal")
	}
}
