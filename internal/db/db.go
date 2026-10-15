package db

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open memilih dialek dari driver. Satu-satunya tempat branch per driver;
// repo/service dilarang branch (lihat ARCHITECTURE.md Bab 2).
func Open(driver, databaseURL string) (*gorm.DB, error) {
	switch strings.ToLower(driver) {
	case "postgres":
		return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	case "sqlite", "":
		dsn := strings.TrimPrefix(databaseURL, "file:")
		if dsn == "" {
			dsn = "dev.db"
		}
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("DB_DRIVER tidak dikenal: %q (sqlite|postgres)", driver)
	}
}
