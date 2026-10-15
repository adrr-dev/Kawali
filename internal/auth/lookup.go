package auth

import (
	"gorm.io/gorm"

	"mewali/internal/models"
)

// Lookup user dipusatkan di sini agar paket lain (users, dsb.) boleh impor
// auth tanpa siklus.

// FindByEmail mencari user dari email (login filler).
func FindByEmail(gdb *gorm.DB, email string) (models.User, error) {
	var u models.User
	err := gdb.Where("email = ?", email).First(&u).Error
	return u, err
}

// FindByID mencari user dari ID (middleware sesi).
func FindByID(gdb *gorm.DB, id uint) (models.User, error) {
	var u models.User
	err := gdb.First(&u, id).Error
	return u, err
}
