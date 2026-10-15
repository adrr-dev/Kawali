package catalog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/media"
	"mewali/internal/models"
	"mewali/internal/orders"
	"mewali/internal/reviews"
	"mewali/internal/sizematch"
	"mewali/internal/users"
	"mewali/web/templates/pages"
)

// CoverFallback dipakai kartu/galeri saat produk belum punya foto.
const CoverFallback = "/static/img/placeholder.svg"

// RegisterRoutes memasang / (beranda), /cari, /produk/:id, dan POST /cart (login).
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	r.GET("/", func(c *gin.Context) {
		baru, _ := Search(gdb, "", "", 6)
		pick, _ := Search(gdb, "", "", 4)
		var banners []models.Banner
		gdb.Where("active = ?", true).Order("sort ASC").Find(&banners)
		d := pages.HomeData{Banners: banners, Baru: baru, Pick: pick, Covers: map[uint]string{}}
		for _, p := range append(baru, pick...) {
			if _, ok := d.Covers[p.ID]; !ok {
				d.Covers[p.ID] = coverURL(gdb, p.ID)
			}
		}
		c.Status(http.StatusOK)
		_ = pages.Home(d).Render(c.Request.Context(), c.Writer)
	})

	r.GET("/cari", func(c *gin.Context) {
		items, err := Search(gdb, c.Query("cat"), c.Query("q"), 24)
		if err != nil {
			c.String(http.StatusInternalServerError, "pencarian gagal")
			return
		}
		d := pages.SearchData{Q: c.Query("q"), Category: c.Query("cat"), Items: items, Covers: map[uint]string{}}
		for _, p := range items {
			d.Covers[p.ID] = coverURL(gdb, p.ID)
		}
		c.Status(http.StatusOK)
		_ = pages.Search(d).Render(c.Request.Context(), c.Writer)
	})

	r.GET("/produk/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var p models.Product
		if err := gdb.First(&p, uint(id)).Error; err != nil {
			c.String(http.StatusNotFound, "produk tidak ditemukan")
			return
		}
		var seller models.User
		gdb.First(&seller, p.SellerID)
		var defects []models.Defect
		gdb.Where("product_id = ?", p.ID).Order("id ASC").Find(&defects)
		var cs models.ConditionScore
		scored := gdb.Where("product_id = ?", p.ID).First(&cs).Error == nil
		avg, _, rated := reviews.SellerRating(gdb, p.SellerID)
		rating := "Belum ada rating"
		if rated {
			rating = comma1(avg) + " ★"
		}
		d := pages.DetailData{Product: p, Seller: seller, Cover: coverURL(gdb, p.ID),
			Defects: defects, Rating: rating}
		if scored {
			d.Scored = true
			d.Score = cs.Score
		}
		if uid, ok := auth.CurrentUserID(c, secret); ok {
			d.Logged = true
			d.IsOwn = seller.ID == uid
			sp := users.GetPreference(gdb, uid)
			ps := GetSizes(gdb, p.ID)
			d.Match = sizematch.Score(
				sizematch.Sizes{Dada: sp.Dada, Pinggang: sp.Pinggang, Pinggul: sp.Pinggul, Tinggi: sp.Tinggi, Torso: sp.Torso},
				sizematch.Sizes{Dada: ps.Dada, Pinggang: ps.Pinggang, Pinggul: ps.Pinggul, Tinggi: ps.Tinggi, Torso: ps.Torso},
			)
		}
		c.Status(http.StatusOK)
		_ = pages.Detail(d).Render(c.Request.Context(), c.Writer)
	})

	r.POST("/cart", auth.RequireLogin(gdb, secret), func(c *gin.Context) {
		u, _ := c.Get("user")
		pid, _ := strconv.ParseUint(c.PostForm("product_id"), 10, 64)
		if err := orders.AddToCart(gdb, u.(models.User).ID, uint(pid)); err != nil {
			c.String(http.StatusBadRequest, "gagal: %v", err)
			return
		}
		c.Redirect(http.StatusFound, "/produk/"+strconv.FormatUint(pid, 10))
	})
}

func coverURL(gdb *gorm.DB, productID uint) string {
	return media.CoverURL(gdb, productID)
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
