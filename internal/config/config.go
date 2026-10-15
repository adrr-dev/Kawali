package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DBDriver      string // sqlite | postgres
	DatabaseURL   string
	SessionSecret string
	UploadDir     string
	AutoMigrate   bool // true = AutoMigrate + seed saat start (dev)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() Config {
	auto, _ := strconv.ParseBool(getenv("AUTO_MIGRATE", "true"))
	return Config{
		Port:          getenv("PORT", "8080"),
		DBDriver:      getenv("DB_DRIVER", "sqlite"),
		DatabaseURL:   getenv("DATABASE_URL", "file:./dev.db"),
		SessionSecret: getenv("SESSION_SECRET", "dev-secret-ganti-di-prod"),
		UploadDir:     getenv("UPLOAD_DIR", "./web/static/uploads"),
		AutoMigrate:   auto,
	}
}
