package users

import (
	"gorm.io/gorm"

	"mewali/internal/models"
)

// GetPreference mengambil preferensi ukuran (bisa kosong/zero bila belum diisi).
func GetPreference(gdb *gorm.DB, userID uint) models.SizePreference {
	var sp models.SizePreference
	gdb.Where("user_id = ?", userID).First(&sp)
	return sp
}

// UpsertPreference menyimpan preferensi ukuran (satu per user).
func UpsertPreference(gdb *gorm.DB, userID uint, sp models.SizePreference) error {
	sp.UserID = userID
	var cur models.SizePreference
	if err := gdb.Where("user_id = ?", userID).First(&cur).Error; err != nil {
		return gdb.Create(&sp).Error
	}
	sp.ID = cur.ID
	return gdb.Save(&sp).Error
}

// UpdateLocation menyimpan label + koordinat manual pengguna.
func UpdateLocation(gdb *gorm.DB, userID uint, label string, lat, lng float64) error {
	return gdb.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{
		"location_label": label, "lat": lat, "lng": lng,
	}).Error
}
