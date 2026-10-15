package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mewali/internal/models"
)

var ErrInvalidLogin = errors.New("email atau sandi salah")

// Login memeriksa email+sandi filler. Selalu kembalikan ErrInvalidLogin bila
// gagal agar tidak membocorkan user mana yang terdaftar.
func Login(gdb *gorm.DB, email, password string) (models.User, error) {
	u, err := FindByEmail(gdb, email)
	if err != nil {
		return models.User{}, ErrInvalidLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return models.User{}, ErrInvalidLogin
	}
	return u, nil
}
