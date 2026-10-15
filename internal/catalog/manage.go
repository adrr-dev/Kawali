package catalog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/models"
	"mewali/internal/storage"
	"mewali/web/templates/pages"
)

// RegisterManage memasang alur tambah/edit katalog (wajib login + pemilik).
func RegisterManage(r *gin.Engine, gdb *gorm.DB, secret string, st storage.Storage) {
	login := auth.RequireLogin(gdb, secret)

	r.GET("/tambah", login, func(c *gin.Context) {
		c.Status(http.StatusOK)
		_ = pages.Tambah("").Render(c.Request.Context(), c.Writer)
	})

	r.POST("/katalog", login, func(c *gin.Context) {
		u := currentSeller(c)
		price, _ := strconv.Atoi(c.PostForm("price"))
		p, err := CreateProduct(gdb, u.ID, c.PostForm("name"), c.PostForm("category"), price)
		if err != nil {
			c.Status(http.StatusBadRequest)
			_ = pages.Tambah(err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}
		gdb.Model(&p).Update("description", c.PostForm("description"))
		_ = SaveSizes(gdb, p.ID, models.ProductSize{
			Dada: parseF(c.PostForm("dada")), Pinggang: parseF(c.PostForm("pinggang")),
			Pinggul: parseF(c.PostForm("pinggul")), Tinggi: parseF(c.PostForm("tinggi")),
			Torso: parseF(c.PostForm("torso")),
		})
		c.Redirect(http.StatusFound, "/katalog/"+strconv.FormatUint(uint64(p.ID), 10))
	})

	r.GET("/katalog/:id", login, func(c *gin.Context) {
		p, ok := ownProduct(c, gdb)
		if !ok {
			return
		}
		renderEdit(c, gdb, p, "")
	})

	r.POST("/katalog/:id/foto", login, func(c *gin.Context) {
		p, ok := ownProduct(c, gdb)
		if !ok {
			return
		}
		slot := c.PostForm("slot")
		fh, err := c.FormFile("photo")
		if err != nil {
			renderEdit(c, gdb, p, "file tidak terbaca")
			return
		}
		f, err := fh.Open()
		if err != nil {
			renderEdit(c, gdb, p, "file tidak terbuka")
			return
		}
		defer f.Close()
		head := make([]byte, 512)
		n, _ := f.Read(head)
		ctype := http.DetectContentType(head[:n])
		if _, err := f.Seek(0, 0); err != nil {
			renderEdit(c, gdb, p, "file tidak terbaca")
			return
		}
		url, err := st.Save(c.Request.Context(), f, ctype)
		if err != nil {
			renderEdit(c, gdb, p, err.Error())
			return
		}
		if err := AddPhoto(gdb, p.ID, slot, url); err != nil {
			renderEdit(c, gdb, p, err.Error())
			return
		}
		c.Redirect(http.StatusFound, "/katalog/"+c.Param("id"))
	})

	r.POST("/katalog/:id/cacat", login, func(c *gin.Context) {
		p, ok := ownProduct(c, gdb)
		if !ok {
			return
		}
		phID, _ := strconv.ParseUint(c.PostForm("photo_id"), 10, 64)
		x, _ := strconv.Atoi(c.PostForm("x"))
		y, _ := strconv.Atoi(c.PostForm("y"))
		var ph models.ProductPhoto
		if err := gdb.Where("id = ? AND product_id = ?", uint(phID), p.ID).First(&ph).Error; err != nil {
			renderEdit(c, gdb, p, "foto tidak milik produk ini")
			return
		}
		if err := AddDefect(gdb, p.ID, ph.ID, x, y, c.PostForm("label")); err != nil {
			renderEdit(c, gdb, p, err.Error())
			return
		}
		c.Redirect(http.StatusFound, "/katalog/"+c.Param("id"))
	})

	r.POST("/katalog/:id/publish", login, func(c *gin.Context) {
		p, ok := ownProduct(c, gdb)
		if !ok {
			return
		}
		if _, err := Publish(gdb, p.ID); err != nil {
			renderEdit(c, gdb, p, err.Error())
			return
		}
		c.Redirect(http.StatusFound, "/produk/"+c.Param("id"))
	})

	r.GET("/barang-dijual", login, func(c *gin.Context) {
		u := currentSeller(c)
		items, _ := SellerList(gdb, u.ID)
		covers := map[uint]string{}
		for _, p := range items {
			covers[p.ID] = coverURL(gdb, p.ID)
		}
		c.Status(http.StatusOK)
		_ = pages.BarangDijual(pages.BarangData{Items: items, Covers: covers}).Render(c.Request.Context(), c.Writer)
	})
}

func currentSeller(c *gin.Context) models.User {
	u, _ := c.Get("user")
	return u.(models.User)
}

// ownProduct memuat produk + memastikan milik penjual yang login.
func ownProduct(c *gin.Context, gdb *gorm.DB) (models.Product, bool) {
	var p models.Product
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := gdb.First(&p, uint(id)).Error; err != nil {
		c.String(http.StatusNotFound, "produk tidak ditemukan")
		return p, false
	}
	if p.SellerID != currentSeller(c).ID {
		c.String(http.StatusForbidden, "bukan barangmu")
		return p, false
	}
	return p, true
}

func renderEdit(c *gin.Context, gdb *gorm.DB, p models.Product, errMsg string) {
	filled, _ := FilledSlots(gdb, p.ID)
	photos, _ := PhotosBySlot(gdb, p.ID)
	var all []models.ProductPhoto
	gdb.Where("product_id = ?", p.ID).Order("id ASC").Find(&all)
	var defects []models.Defect
	gdb.Where("product_id = ?", p.ID).Order("id ASC").Find(&defects)
	d := pages.EditData{Product: p, Slots: RequiredSlots, Filled: filled,
		Missing: MissingSlots(filled), Photos: photos, AllPics: all,
		Defects: defects, Sizes: GetSizes(gdb, p.ID), Error: errMsg}
	c.Status(http.StatusOK)
	_ = pages.EditKatalog(d).Render(c.Request.Context(), c.Writer)
}

func parseF(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	if f < 0 {
		return 0
	}
	return f
}
