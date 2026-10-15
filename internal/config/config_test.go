package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_DRIVER", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("SESSION_SECRET", "")
	t.Setenv("UPLOAD_DIR", "")
	t.Setenv("AUTO_MIGRATE", "")

	cfg := Load()

	cases := map[string][2]string{
		"Port":          {cfg.Port, "8080"},
		"DBDriver":      {cfg.DBDriver, "sqlite"},
		"DatabaseURL":   {cfg.DatabaseURL, "file:./dev.db"},
		"SessionSecret": {cfg.SessionSecret, "dev-secret-ganti-di-prod"},
		"UploadDir":     {cfg.UploadDir, "./web/static/uploads"},
	}
	if !cfg.AutoMigrate {
		t.Error("AutoMigrate default harus true")
	}
	for name, c := range cases {
		if c[0] != c[1] {
			t.Errorf("%s = %q, mau %q", name, c[0], c[1])
		}
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("DB_DRIVER", "postgres")
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/x?sslmode=disable")
	t.Setenv("SESSION_SECRET", "rahasia")
	t.Setenv("UPLOAD_DIR", "/data/up")

	cfg := Load()

	if cfg.Port != "9000" {
		t.Errorf("Port = %q", cfg.Port)
	}
	if cfg.DBDriver != "postgres" {
		t.Errorf("DBDriver = %q", cfg.DBDriver)
	}
	if cfg.DatabaseURL != "postgres://localhost:5432/x?sslmode=disable" {
		t.Errorf("DatabaseURL = %q", cfg.DatabaseURL)
	}
	if cfg.SessionSecret != "rahasia" {
		t.Errorf("SessionSecret = %q", cfg.SessionSecret)
	}
	if cfg.UploadDir != "/data/up" {
		t.Errorf("UploadDir = %q", cfg.UploadDir)
	}
}

func TestLoadPartialEnvFallsBackToDefault(t *testing.T) {
	t.Setenv("PORT", "9001")
	t.Setenv("DB_DRIVER", "")

	cfg := Load()

	if cfg.Port != "9001" {
		t.Errorf("Port = %q", cfg.Port)
	}
	if cfg.DBDriver != "sqlite" {
		t.Errorf("DBDriver = %q, mau default sqlite", cfg.DBDriver)
	}
}
