package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveRoundtrip(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "up")
	l := Local{Dir: dir, BaseURL: "/static/uploads"}
	url, err := l.Save(context.Background(), bytes.NewReader([]byte("fake-jpeg-data")), "image/jpeg")
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if !strings.HasPrefix(url, "/static/uploads/") || !strings.HasSuffix(url, ".jpg") {
		t.Errorf("url aneh: %q", url)
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.Base(url))); err != nil {
		t.Errorf("file tidak ada: %v", err)
	}
}

func TestSaveRejects(t *testing.T) {
	l := Local{Dir: t.TempDir(), BaseURL: "/u"}
	if _, err := l.Save(context.Background(), bytes.NewReader([]byte("x")), "image/gif"); err == nil {
		t.Error("gif harus ditolak")
	}
	big := make([]byte, MaxFile+10)
	if _, err := l.Save(context.Background(), bytes.NewReader(big), "image/png"); err == nil {
		t.Error("file raksasa harus ditolak")
	}
}
