package reports

import (
	"errors"

	"gorm.io/gorm"

	"mewali/internal/models"
)

var (
	ErrNoTarget = errors.New("tentukan pengguna atau produk yang dilaporkan")
	ErrNoReason = errors.New("alasan laporan wajib diisi")
)

// FileReport mengirim laporan untuk ditinjau (status awal open).
// Satu dari targetUserID / productID wajib diisi.
func FileReport(gdb *gorm.DB, reporterID uint, targetUserID, productID *uint, reason, detail string) error {
	if targetUserID == nil && productID == nil {
		return ErrNoTarget
	}
	if reason == "" {
		return ErrNoReason
	}
	return gdb.Create(&models.Report{ReporterID: reporterID, TargetUserID: targetUserID,
		ProductID: productID, Reason: reason, Detail: detail, Status: models.ReportOpen}).Error
}
