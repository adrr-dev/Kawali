package geo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/models"
	"mewali/web/templates/pages"
)

// RegisterRoutes memasang /peta: daftar produk + titik donasi dengan jarak.
// Peta adalah permukaan statis pucat (tanpa vendor/GPS), fallback daftar.
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	login := auth.RequireLogin(gdb, secret)

	r.GET("/peta", login, func(c *gin.Context) {
		u := currentViewer(c)
		lat, lng := locate(c, u)
		var products []models.Product
		gdb.Where("status = ?", models.ProductActive).Order("id DESC").Limit(30).Find(&products)
		var donations []models.DonationPoint
		gdb.Where("active = ?", true).Find(&donations)
		d := pages.PetaData{Lat: lat, Lng: lng, HasLoc: lat != 0 || lng != 0}
		for _, p := range products {
			d.Pins = append(d.Pins, pages.PetaPin{Name: p.Name, URL: "/produk/" + strconv.FormatUint(uint64(p.ID), 10),
				Kind: "produk", Dist: FormatID(HaversineKm(lat, lng, p.Lat, p.Lng))})
		}
		for _, dn := range donations {
			d.Pins = append(d.Pins, pages.PetaPin{Name: dn.Name, URL: "#", Kind: "donasi",
				Dist: FormatID(HaversineKm(lat, lng, dn.Lat, dn.Lng))})
		}
		c.Status(http.StatusOK)
		_ = pages.Peta(d).Render(c.Request.Context(), c.Writer)
	})
}

func currentViewer(c *gin.Context) models.User {
	u, _ := c.Get("user")
	return u.(models.User)
}

// locate memakai lokasi user, dioverride query lat/lng manual bila ada.
func locate(c *gin.Context, u models.User) (float64, float64) {
	lat, lng := u.Lat, u.Lng
	if s := c.Query("lat"); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			lat = f
		}
	}
	if s := c.Query("lng"); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			lng = f
		}
	}
	return lat, lng
}
