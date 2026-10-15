package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mewali/internal/db"
	"mewali/internal/models"
)

func testEngine(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "h.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	gdb.Create(&models.User{Name: "T", Email: "t@test.id", Handle: "t_1", PasswordHash: string(hash)})

	r := gin.New()
	RegisterRoutes(r, gdb, "test-secret")
	r.GET("/khusus", RequireLogin(gdb, "test-secret"), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r, gdb
}

func loginCookie(t *testing.T, r *gin.Engine, email, pass string) *http.Cookie {
	t.Helper()
	form := url.Values{"email": {email}, "password": {pass}}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("login %s: status = %d, mau 302", email, w.Code)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == CookieName {
			return c
		}
	}
	t.Fatal("cookie sesi tidak diset")
	return nil
}

func TestLoginFlowSetsCookie(t *testing.T) {
	r, _ := testEngine(t)
	c := loginCookie(t, r, "t@test.id", "password")
	if c.Value == "" {
		t.Error("cookie kosong")
	}
}

func TestLoginWrongPassword401(t *testing.T) {
	r, _ := testEngine(t)
	form := url.Values{"email": {"t@test.id"}, "password": {"salah"}}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, mau 401", w.Code)
	}
}

func TestRequireLoginGate(t *testing.T) {
	r, _ := testEngine(t)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/khusus", nil))
	if w.Code != http.StatusFound {
		t.Errorf("anonim: status = %d, mau 302", w.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/khusus", nil)
	req.AddCookie(loginCookie(t, r, "t@test.id", "password"))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	if w2.Code != http.StatusOK {
		t.Errorf("login: status = %d, mau 200", w2.Code)
	}
}
