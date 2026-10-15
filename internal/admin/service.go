package admin

import (
	"errors"

	"gorm.io/gorm"

	"mewali/internal/models"
)

// Overview untuk ringkasan operasional (hanya dari data yang ada).
type Overview struct {
	Users        int64
	Products     int64
	PendingOrder int64
	DoneOrder    int64
	OpenReports  int64
}

// GetOverview menghitung metrik dasar platform.
func GetOverview(gdb *gorm.DB) Overview {
	var o Overview
	gdb.Model(&models.User{}).Count(&o.Users)
	gdb.Model(&models.Product{}).Count(&o.Products)
	gdb.Model(&models.Order{}).Where("status = ?", models.OrderPending).Count(&o.PendingOrder)
	gdb.Model(&models.Order{}).Where("status = ?", models.OrderDone).Count(&o.DoneOrder)
	gdb.Model(&models.Report{}).Where("status IN ?", []string{models.ReportOpen, models.ReportInReview}).Count(&o.OpenReports)
	return o
}

// ResolveReport menutup laporan (resolve=true) atau menolaknya.
// Suspend/blokir akun sengaja TIDAK ada (tunggu kebijakan produk).
func ResolveReport(gdb *gorm.DB, adminID, reportID uint, resolve bool) error {
	status := models.ReportRejected
	action := "report_reject"
	if resolve {
		status = models.ReportResolved
		action = "report_resolve"
	}
	return gdb.Transaction(func(tx *gorm.DB) error {
		r := tx.Model(&models.Report{}).Where("id = ?", reportID).Update("status", status)
		if r.Error != nil {
			return r.Error
		}
		if r.RowsAffected == 0 {
			return errors.New("laporan tidak ditemukan")
		}
		return audit(tx, adminID, action, "report", reportID, "")
	})
}

// Takedown menurunkan listing (taken_down). Bukan hapus permanen.
func Takedown(gdb *gorm.DB, adminID, productID uint) error {
	return gdb.Transaction(func(tx *gorm.DB) error {
		r := tx.Model(&models.Product{}).Where("id = ?", productID).Update("status", models.ProductTakenDown)
		if r.Error != nil {
			return r.Error
		}
		if r.RowsAffected == 0 {
			return errors.New("produk tidak ditemukan")
		}
		return audit(tx, adminID, "product_takedown", "product", productID, "")
	})
}

func audit(tx *gorm.DB, adminID uint, action, targetType string, targetID uint, note string) error {
	return tx.Create(&models.AdminAudit{AdminID: adminID, Action: action, TargetType: targetType, TargetID: targetID, Note: note}).Error
}
