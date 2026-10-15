package users

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mewali/internal/models"
)

var (
	ErrEmailTaken = errors.New("email sudah terdaftar")
	ErrBadSignup  = errors.New("nama, email, dan sandi (min. 8) wajib diisi")
)

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// Register membuat akun baru + handle unik turunan nama.
func Register(gdb *gorm.DB, name, email, password string) (models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	if name == "" || email == "" || len(password) < 8 {
		return models.User{}, ErrBadSignup
	}
	var n int64
	gdb.Model(&models.User{}).Where("email = ?", email).Count(&n)
	if n > 0 {
		return models.User{}, ErrEmailTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	base := nonAlnum.ReplaceAllString(strings.ToLower(name), "")
	if base == "" {
		base = "pengguna"
	}
	handle := base
	for i := 2; ; i++ {
		var c int64
		gdb.Model(&models.User{}).Where("handle = ?", handle).Count(&c)
		if c == 0 {
			break
		}
		handle = base + strconv.Itoa(i)
	}
	u := models.User{Name: name, Email: email, PasswordHash: string(hash), Handle: handle, Level: "Pemula Hijau"}
	return u, gdb.Create(&u).Error
}
