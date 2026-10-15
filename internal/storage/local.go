package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MaxFile 5MB sesuai ARCHITECTURE Bab 3.
const MaxFile = 5 << 20

var allowed = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// Storage menyimpan upload dan mengembalikan URL publik. Skema DB tetap
// string URL sehingga ganti backend (Supabase) tanpa migrasi.
type Storage interface {
	Save(ctx context.Context, r io.Reader, contentType string) (url string, err error)
}

// Local menyimpan ke direktori disk yang diserve statis.
type Local struct {
	Dir     string // mis. ./web/static/uploads
	BaseURL string // mis. /static/uploads
}

func (l Local) Save(ctx context.Context, r io.Reader, contentType string) (string, error) {
	ext, ok := allowed[strings.ToLower(contentType)]
	if !ok {
		return "", errors.New("tipe file harus jpeg/png/webp")
	}
	var buf [MaxFile + 1]byte
	n, err := io.ReadFull(r, buf[:])
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}
	if n > MaxFile {
		return "", errors.New("file maksimal 5MB")
	}
	var name [8]byte
	if _, err := rand.Read(name[:]); err != nil {
		return "", err
	}
	filename := hex.EncodeToString(name[:]) + ext
	if err := os.MkdirAll(l.Dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(l.Dir, filename), buf[:n], 0o644); err != nil {
		return "", err
	}
	return l.BaseURL + "/" + filename, nil
}
