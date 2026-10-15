package main

import (
	"log"
	"net/http"

	"mewali/internal/admin"
	"mewali/internal/auth"
	"mewali/internal/catalog"
	"mewali/internal/chat"
	"mewali/internal/config"
	"mewali/internal/db"
	"mewali/internal/geo"
	"mewali/internal/orders"
	"mewali/internal/points"
	"mewali/internal/seed"
	"mewali/internal/storage"
	"mewali/internal/users"
	"mewali/web/templates/pages"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gdb, err := db.Open(cfg.DBDriver, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("buka db: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("sql db: %v", err)
	}
	defer sqlDB.Close()

	if cfg.AutoMigrate {
		if err := db.Migrate(gdb); err != nil {
			log.Fatalf("migrate: %v", err)
		}
		if err := seed.Run(gdb); err != nil {
			log.Fatalf("seed: %v", err)
		}
	}

	r := gin.Default()
	r.Static("/static", "./web/static")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true, "app": "mewali"})
	})
	r.GET("/onboarding", func(c *gin.Context) {
		c.Status(http.StatusOK)
		_ = pages.Onboarding().Render(c.Request.Context(), c.Writer)
	})
	auth.RegisterRoutes(r, gdb, cfg.SessionSecret)
	catalog.RegisterRoutes(r, gdb, cfg.SessionSecret)
	st := storage.Local{Dir: cfg.UploadDir, BaseURL: "/static/uploads"}
	catalog.RegisterManage(r, gdb, cfg.SessionSecret, st)
	orders.RegisterRoutes(r, gdb, cfg.SessionSecret)
	chat.RegisterRoutes(r, gdb, cfg.SessionSecret)
	points.RegisterRoutes(r, gdb, cfg.SessionSecret)
	users.RegisterRoutes(r, gdb, cfg.SessionSecret)
	geo.RegisterRoutes(r, gdb, cfg.SessionSecret)
	admin.RegisterRoutes(r, gdb, cfg.SessionSecret)

	log.Printf("listen :%s (db=%s)", cfg.Port, cfg.DBDriver)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
