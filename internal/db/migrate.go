package db

import (
	"gorm.io/gorm"

	"mewali/internal/models"
)

// Migrate menjalankan AutoMigrate untuk dev/seed. Prod memakai
// migrations/*.sql berversion (lihat ARCHITECTURE.md Bab 2).
func Migrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(models.All()...)
}
