package users

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/media"
	"mewali/internal/models"
	"mewali/internal/points"
	"mewali/internal/reports"
	"mewali/internal/reviews"
	"mewali/web/templates/pages"
)

// RegisterRoutes memasang profil, preferensi, ulasan, dan laporan.
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	login := auth.RequireLogin(gdb, secret)

	r.GET("/profil", login, func(c *gin.Context) {
		u := fresh(c, gdb)
		items, _ := ownProducts(gdb, u.ID)
		names := reviewerNames(gdb, u.ID)
		var revs []models.Review
		gdb.Where("seller_id = ?", u.ID).Order("id DESC").Limit(5).Find(&revs)
		d := pages.ProfilData{User: u, ProductCount: productCount(gdb, u.ID),
			Rating: ratingOf(gdb, u.ID), Balance: points.Balance(gdb, u.ID),
			Level: u.Level, Covers: coversOf(gdb, items), Reviews: reviewViews(revs, names)}
		for _, p := range items {
			if len(d.Items) >= 4 {
				break
			}
			d.Items = append(d.Items, p)
		}
		c.Status(http.StatusOK)
		_ = pages.Profil(d).Render(c.Request.Context(), c.Writer)
	})

	r.GET("/profil/:handle", login, func(c *gin.Context) {
		me := fresh(c, gdb)
		target, err := FindByHandle(gdb, c.Param("handle"))
		if err != nil {
			c.String(http.StatusNotFound, "pengguna tidak ditemukan")
			return
		}
		if target.ID == me.ID {
			c.Redirect(http.StatusFound, "/profil")
			return
		}
		var items []models.Product
		gdb.Where("seller_id = ? AND status = ?", target.ID, models.ProductActive).Order("id DESC").Limit(8).Find(&items)
		var revs []models.Review
		gdb.Where("seller_id = ?", target.ID).Order("id DESC").Limit(10).Find(&revs)
		d := pages.ProfilOrangData{User: target, Rating: ratingOf(gdb, target.ID),
			Items: items, Covers: coversOf(gdb, items),
			Reviews:  reviewViews(revs, reviewerNames(gdb, target.ID)),
			Eligible: eligibleOrders(gdb, me.ID, target.ID)}
		c.Status(http.StatusOK)
		_ = pages.ProfilOrang(d).Render(c.Request.Context(), c.Writer)
	})

	r.GET("/preferensi", login, func(c *gin.Context) {
		u := fresh(c, gdb)
		sp := GetPreference(gdb, u.ID)
		d := pages.PreferensiData{Sizes: sp, Label: u.LocationLabel, Lat: u.Lat, Lng: u.Lng}
		c.Status(http.StatusOK)
		_ = pages.Preferensi(d).Render(c.Request.Context(), c.Writer)
	})

	r.POST("/preferensi", login, func(c *gin.Context) {
		u := fresh(c, gdb)
		_ = UpsertPreference(gdb, u.ID, models.SizePreference{
			Dada: parseF(c.PostForm("dada")), Pinggang: parseF(c.PostForm("pinggang")),
			Pinggul: parseF(c.PostForm("pinggul")), Tinggi: parseF(c.PostForm("tinggi")),
			Torso: parseF(c.PostForm("torso")),
		})
		lat, _ := strconv.ParseFloat(c.PostForm("lat"), 64)
		lng, _ := strconv.ParseFloat(c.PostForm("lng"), 64)
		_ = UpdateLocation(gdb, u.ID, c.PostForm("location"), lat, lng)
		c.Redirect(http.StatusFound, "/profil")
	})

	r.POST("/ulasan", login, func(c *gin.Context) {
		u := fresh(c, gdb)
		oid, _ := strconv.ParseUint(c.PostForm("order_id"), 10, 64)
		stars, _ := strconv.Atoi(c.PostForm("stars"))
		if err := reviews.CreateReview(gdb, uint(oid), u.ID, stars, c.PostForm("body")); err != nil {
			c.String(http.StatusBadRequest, "gagal: %v", err)
			return
		}
		var o models.Order
		gdb.First(&o, uint(oid))
		var seller models.User
		gdb.First(&seller, o.SellerID)
		c.Redirect(http.StatusFound, "/profil/"+seller.Handle)
	})

	r.POST("/lapor", login, func(c *gin.Context) {
		u := fresh(c, gdb)
		var tuid, pid *uint
		if s := c.PostForm("target_user_id"); s != "" {
			if n, err := strconv.ParseUint(s, 10, 64); err == nil {
				v := uint(n)
				tuid = &v
			}
		}
		if s := c.PostForm("product_id"); s != "" {
			if n, err := strconv.ParseUint(s, 10, 64); err == nil {
				v := uint(n)
				pid = &v
			}
		}
		back := c.PostForm("back")
		if back == "" {
			back = "/"
		}
		if err := reports.FileReport(gdb, u.ID, tuid, pid, c.PostForm("reason"), c.PostForm("detail")); err != nil {
			c.String(http.StatusBadRequest, "gagal: %v", err)
			return
		}
		c.Redirect(http.StatusFound, back)
	})

	r.GET("/signup", func(c *gin.Context) {
		c.Status(http.StatusOK)
		_ = pages.Signup("").Render(c.Request.Context(), c.Writer)
	})

	r.POST("/signup", func(c *gin.Context) {
		u, err := Register(gdb, c.PostForm("name"), c.PostForm("email"), c.PostForm("password"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			_ = pages.Signup(err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}
		val, age := auth.NewSession(secret, u.ID, time.Now())
		c.SetCookie(auth.CookieName, val, age, "/", "", false, true)
		c.Redirect(http.StatusFound, "/preferensi")
	})

	r.GET("/lupa-sandi", func(c *gin.Context) {
		c.Status(http.StatusOK)
		_ = pages.Lupa(false).Render(c.Request.Context(), c.Writer)
	})

	r.POST("/lupa-sandi", func(c *gin.Context) {
		// Selalu respons netral: tidak membocorkan email terdaftar atau tidak.
		c.Status(http.StatusOK)
		_ = pages.Lupa(true).Render(c.Request.Context(), c.Writer)
	})
}

func fresh(c *gin.Context, gdb *gorm.DB) models.User {
	u, _ := c.Get("user")
	var f models.User
	gdb.First(&f, u.(models.User).ID)
	return f
}

// FindByHandle mencari user dari handle publik.
func FindByHandle(gdb *gorm.DB, handle string) (models.User, error) {
	var u models.User
	err := gdb.Where("handle = ?", handle).First(&u).Error
	return u, err
}

func ownProducts(gdb *gorm.DB, sellerID uint) ([]models.Product, error) {
	var out []models.Product
	err := gdb.Where("seller_id = ?", sellerID).Order("id DESC").Find(&out).Error
	return out, err
}

func productCount(gdb *gorm.DB, sellerID uint) int64 {
	var n int64
	gdb.Model(&models.Product{}).Where("seller_id = ? AND status = ?", sellerID, models.ProductActive).Count(&n)
	return n
}

func ratingOf(gdb *gorm.DB, sellerID uint) string {
	avg, _, rated := reviews.SellerRating(gdb, sellerID)
	if !rated {
		return "Belum ada rating"
	}
	return comma1(avg) + " ★"
}

func comma1(f float64) string {
	s := strconv.FormatFloat(f, 'f', 1, 64)
	out := ""
	for _, ch := range s {
		if ch == '.' {
			out += ","
		} else {
			out += string(ch)
		}
	}
	return out
}

func coversOf(gdb *gorm.DB, items []models.Product) map[uint]string {
	m := map[uint]string{}
	for _, p := range items {
		m[p.ID] = media.CoverURL(gdb, p.ID)
	}
	return m
}

func reviewerNames(gdb *gorm.DB, sellerID uint) map[uint]string {
	var revs []models.Review
	gdb.Where("seller_id = ?", sellerID).Select("reviewer_id").Find(&revs)
	m := map[uint]string{}
	for _, r := range revs {
		if _, ok := m[r.ReviewerID]; ok {
			continue
		}
		var u models.User
		if err := gdb.Select("name").First(&u, r.ReviewerID).Error; err == nil {
			m[r.ReviewerID] = u.Name
		}
	}
	return m
}

func reviewViews(revs []models.Review, names map[uint]string) []pages.ReviewView {
	var out []pages.ReviewView
	for _, r := range revs {
		out = append(out, pages.ReviewView{ReviewerName: names[r.ReviewerID], Stars: r.Stars, Body: r.Body})
	}
	return out
}

// eligibleOrders: order done milikku ke penjual ini yang belum diulas.
func eligibleOrders(gdb *gorm.DB, buyerID, sellerID uint) []pages.EligibleOrder {
	var done []models.Order
	gdb.Where("buyer_id = ? AND seller_id = ? AND status = ?", buyerID, sellerID, models.OrderDone).Find(&done)
	var out []pages.EligibleOrder
	for _, o := range done {
		var n int64
		gdb.Model(&models.Review{}).Where("order_id = ?", o.ID).Count(&n)
		if n > 0 {
			continue
		}
		var p models.Product
		gdb.Select("name").First(&p, o.ProductID)
		out = append(out, pages.EligibleOrder{OrderID: o.ID, ProductName: p.Name})
	}
	return out
}

func parseF(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	if f < 0 {
		return 0
	}
	return f
}
