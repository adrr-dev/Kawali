package orders

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/chat"
	"mewali/internal/media"
	"mewali/internal/models"
	"mewali/web/templates/pages"
)

// RegisterRoutes memasang keranjang, beli, checkout, konfirmasi, sukses.
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	login := auth.RequireLogin(gdb, secret)

	r.GET("/keranjang", login, func(c *gin.Context) {
		u := me(c)
		var items []models.CartItem
		gdb.Where("user_id = ?", u.ID).Find(&items)
		d := pages.CartData{Covers: map[uint]string{}}
		for _, it := range items {
			var p models.Product
			if err := gdb.Where("id = ? AND status = ?", it.ProductID, models.ProductActive).First(&p).Error; err != nil {
				continue
			}
			d.Items = append(d.Items, p)
			d.Covers[p.ID] = coverOf(gdb, p.ID)
			d.Total += p.Price
		}
		c.Status(http.StatusOK)
		_ = pages.Keranjang(d).Render(c.Request.Context(), c.Writer)
	})

	r.POST("/beli", login, func(c *gin.Context) {
		u := me(c)
		pid, _ := strconv.ParseUint(c.PostForm("product_id"), 10, 64)
		var p models.Product
		if err := gdb.First(&p, uint(pid)).Error; err != nil {
			c.String(http.StatusNotFound, "produk tidak ditemukan")
			return
		}
		renderBeli(c, gdb, u, p, "")
	})

	r.POST("/checkout", login, func(c *gin.Context) {
		u := me(c)
		pid, _ := strconv.ParseUint(c.PostForm("product_id"), 10, 64)
		var p models.Product
		if err := gdb.First(&p, uint(pid)).Error; err != nil {
			c.String(http.StatusNotFound, "produk tidak ditemukan")
			return
		}
		if u.Lat == 0 && u.Lng == 0 {
			renderBeli(c, gdb, u, p, "Atur lokasi dulu sebelum membayar.")
			return
		}
		o, _, _, err := Checkout(gdb, u.ID, p.ID, c.PostForm("pay"), c.PostForm("mode"), time.Now())
		if err != nil {
			renderBeli(c, gdb, u, p, err.Error())
			return
		}
		if o.Mode == models.ModeMeetup {
			c.Redirect(http.StatusFound, "/konfirmasi/"+strconv.FormatUint(uint64(o.ID), 10))
			return
		}
		c.Redirect(http.StatusFound, "/transaksi/"+strconv.FormatUint(uint64(o.ID), 10))
	})

	r.GET("/transaksi/:id", login, func(c *gin.Context) {
		o, p, ok := memberOrder(c, gdb)
		if !ok {
			return
		}
		u := me(c)
		th, err := chat.GetOrCreateThread(gdb, p.ID, o.BuyerID, o.SellerID, &o.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "thread gagal")
			return
		}
		d := pages.TransaksiData{Order: o, Product: p, IsBuyer: u.ID == o.BuyerID, ThreadID: th.ID}
		c.Status(http.StatusOK)
		_ = pages.Transaksi(d).Render(c.Request.Context(), c.Writer)
	})

	r.POST("/transaksi/:id/selesai", login, func(c *gin.Context) {
		o, _, ok := memberOrder(c, gdb)
		if !ok {
			return
		}
		if err := P2PDone(gdb, o.ID); err != nil {
			c.String(http.StatusBadRequest, "gagal: %v", err)
			return
		}
		c.Redirect(http.StatusFound, "/sukses/"+strconv.FormatUint(uint64(o.ID), 10))
	})

	r.GET("/konfirmasi/:id", login, func(c *gin.Context) {
		o, _, ok := memberOrder(c, gdb)
		if !ok {
			return
		}
		renderConfirm(c, gdb, secret, o, "")
	})

	r.POST("/konfirmasi/:id", login, func(c *gin.Context) {
		o, _, ok := memberOrder(c, gdb)
		if !ok {
			return
		}
		u := me(c)
		role := models.CodeSeller
		if u.ID == o.BuyerID {
			role = models.CodeBuyer
		}
		if err := Confirm(gdb, o.ID, role, c.PostForm("code"), time.Now()); err != nil {
			renderConfirm(c, gdb, secret, o, confirmMsg(err))
			return
		}
		c.Redirect(http.StatusFound, "/sukses/"+strconv.FormatUint(uint64(o.ID), 10))
	})

	r.GET("/sukses/:id", login, func(c *gin.Context) {
		o, _, ok := memberOrder(c, gdb)
		if !ok {
			return
		}
		pts := 0
		if o.Status == models.OrderDone {
			var led models.PointLedger
			if err := gdb.Where("order_id = ? AND user_id = ?", o.ID, o.BuyerID).First(&led).Error; err == nil {
				pts = led.Delta
			}
		}
		c.Status(http.StatusOK)
		_ = pages.Sukses(pages.SuksesData{Order: o, Points: pts}).Render(c.Request.Context(), c.Writer)
	})
}

func me(c *gin.Context) models.User {
	u, _ := c.Get("user")
	return u.(models.User)
}

// memberOrder memuat order + memastikan user adalah buyer/seller.
func memberOrder(c *gin.Context, gdb *gorm.DB) (models.Order, models.Product, bool) {
	var o models.Order
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := gdb.First(&o, uint(id)).Error; err != nil {
		c.String(http.StatusNotFound, "transaksi tidak ditemukan")
		return o, models.Product{}, false
	}
	u := me(c)
	if u.ID != o.BuyerID && u.ID != o.SellerID {
		c.String(http.StatusForbidden, "bukan transaksimu")
		return o, models.Product{}, false
	}
	var p models.Product
	gdb.First(&p, o.ProductID)
	return o, p, true
}

func renderBeli(c *gin.Context, gdb *gorm.DB, u models.User, p models.Product, errMsg string) {
	d := pages.BeliData{Product: p, Cover: coverOf(gdb, p.ID), NeedLoc: u.Lat == 0 && u.Lng == 0, Error: errMsg}
	c.Status(http.StatusOK)
	_ = pages.Beli(d).Render(c.Request.Context(), c.Writer)
}

func renderConfirm(c *gin.Context, gdb *gorm.DB, secret string, o models.Order, errMsg string) {
	if o.Mode != models.ModeMeetup || o.Status != models.OrderPending {
		c.Redirect(http.StatusFound, "/transaksi/"+strconv.FormatUint(uint64(o.ID), 10))
		return
	}
	u := me(c)
	ownRole := models.CodeSeller
	buyer := false
	if u.ID == o.BuyerID {
		ownRole = models.CodeBuyer
		buyer = true
	}
	var row models.ConfirmCode
	if err := gdb.Where("order_id = ? AND role = ?", o.ID, ownRole).First(&row).Error; err != nil {
		c.String(http.StatusInternalServerError, "kode belum terbit")
		return
	}
	th, err := chat.GetOrCreateThread(gdb, o.ProductID, o.BuyerID, o.SellerID, &o.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "thread gagal")
		return
	}
	d := pages.KonfirmasiData{OrderID: o.ID, OwnCode: row.Code,
		ExpiresAt: row.ExpiresAt.Unix(), IAmBuyer: buyer, ThreadID: th.ID, Error: errMsg}
	c.Status(http.StatusOK)
	_ = pages.Konfirmasi(d).Render(c.Request.Context(), c.Writer)
}

func confirmMsg(err error) string {
	switch err {
	case ErrCodeWrong:
		return "Kode salah. Minta kode terbaru dari lawan transaksimu."
	case ErrCodeExpired:
		return "Waktu habis. Mulai ulang transaksi bila perlu."
	case ErrCodeUsed:
		return "Transaksi ini sudah selesai."
	default:
		return "Gagal: " + err.Error()
	}
}

func coverOf(gdb *gorm.DB, productID uint) string {
	return media.CoverURL(gdb, productID)
}
