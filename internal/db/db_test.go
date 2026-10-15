package db

import (
	"path/filepath"
	"testing"
)

func TestOpenUnknownDriverFails(t *testing.T) {
	if _, err := Open("mysql", "x"); err == nil {
		t.Error("mau error untuk driver tak dikenal, dapat nil")
	}
}

func TestOpenSQLiteFilePrefix(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	gdb, err := Open("sqlite", "file:"+path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("DB(): %v", err)
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestOpenSQLiteBarePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bare.db")

	gdb, err := Open("sqlite", path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("DB(): %v", err)
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestOpenSQLiteDefaultDriver(t *testing.T) {
	path := filepath.Join(t.TempDir(), "default.db")

	// driver "" harus jatuh ke sqlite (lihat Load: default sqlite)
	gdb, err := Open("", "file:"+path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("DB(): %v", err)
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}
