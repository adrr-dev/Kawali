package chat

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
)

func openTest(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "chat.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return gdb
}

func TestThreadReuseAndMessageOrder(t *testing.T) {
	gdb := openTest(t)
	a, err := GetOrCreateThread(gdb, 9, 2, 1, nil)
	if err != nil {
		t.Fatalf("thread: %v", err)
	}
	b, err := GetOrCreateThread(gdb, 9, 2, 1, nil)
	if err != nil || a.ID != b.ID {
		t.Fatalf("thread harus dipakai ulang: %v", err)
	}
	for _, m := range []struct {
		sender uint
		body   string
	}{{2, "halo"}, {1, "masih ada?"}, {2, "deal"}} {
		if err := PostMessage(gdb, a.ID, m.sender, m.body); err != nil {
			t.Fatalf("post: %v", err)
		}
	}
	msgs, err := ListMessages(gdb, a.ID, 50)
	if err != nil || len(msgs) != 3 {
		t.Fatalf("list: %v (%d)", err, len(msgs))
	}
	if msgs[0].Body != "halo" || msgs[2].Body != "deal" {
		t.Error("urutan harus tertua-dulu")
	}
}

func TestPostGuards(t *testing.T) {
	gdb := openTest(t)
	th, _ := GetOrCreateThread(gdb, 9, 2, 1, nil)
	if err := PostMessage(gdb, th.ID, 99, "x"); err != ErrNotMember {
		t.Errorf("orang luar mau ErrNotMember, dapat %v", err)
	}
	if err := PostMessage(gdb, th.ID, 2, ""); err != ErrEmptyBody {
		t.Errorf("kosong mau ErrEmptyBody, dapat %v", err)
	}
}
