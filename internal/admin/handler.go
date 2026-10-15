package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/models"
	"mewali/web/templates/pages"
)

// RequireAdmin menolak non-admin (403). Tanpa role lain sampai kebijakan ada.
func RequireAdmin(gdb *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := auth.CurrentUserID(c, secret)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		u, err := auth.FindByID(gdb, uid)
		if err != nil || !u.IsAdmin {
			c.String(http.StatusForbidden, "khusus admin")
			c.Abort()
			return
		}
		c.Set("admin", u)
		c.Next()
	}
}

func me(c *gin.Context) models.User {
	u, _ := c.Get("admin")
	return u.(models.User)
}

// RegisterRoutes memasang dasbor admin minimum (baca + moderasi dasar).
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	adm := r.Group("/admin", RequireAdmin(gdb, secret))

	adm.GET("", func(c *gin.Context) {
		c.Status(http.StatusOK)
		_ = pages.AdminOverview(pageOv(GetOverview(gdb))).Render(c.Request.Context(), c.Writer)
	})

	adm.GET("/laporan", func(c *gin.Context) {
		var reps []models.Report
		gdb.Order("id DESC").Limit(50).Find(&reps)
		var rows []pages.AdminReportRow
		for _, x := range reps {
			rows = append(rows, pages.AdminReportRow{ID: x.ID, Reason: x.Reason, Status: x.Status})
		}
		c.Status(http.StatusOK)
		_ = pages.AdminLaporan(rows).Render(c.Request.Context(), c.Writer)
	})

	adm.POST("/laporan/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err := ResolveReport(gdb, me(c).ID, uint(id), c.PostForm("aksi") == "resolve")
		if err != nil {
			c.String(http.StatusBadRequest, "gagal: %v", err)
			return
		}
		c.Redirect(http.StatusFound, "/admin/laporan")
	})

	adm.GET("/produk", func(c *gin.Context) {
		var items []models.Product
		gdb.Order("id DESC").Limit(50).Find(&items)
		var rows []pages.AdminProductRow
		for _, p := range items {
			rows = append(rows, pages.AdminProductRow{ID: p.ID, Name: p.Name, Status: p.Status})
		}
		c.Status(http.StatusOK)
		_ = pages.AdminProduk(rows).Render(c.Request.Context(), c.Writer)
	})

	adm.POST("/produk/:id/takedown", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		if err := Takedown(gdb, me(c).ID, uint(id)); err != nil {
			c.String(http.StatusBadRequest, "gagal: %v", err)
			return
		}
		c.Redirect(http.StatusFound, "/admin/produk")
	})

	adm.GET("/transaksi", func(c *gin.Context) {
		var orders []models.Order
		gdb.Order("id DESC").Limit(50).Find(&orders)
		var rows []pages.AdminOrderRow
		for _, o := range orders {
			rows = append(rows, pages.AdminOrderRow{ID: o.ID, Total: o.Total, Pay: o.PayMethod, Mode: o.Mode, Status: o.Status})
		}
		c.Status(http.StatusOK)
		_ = pages.AdminTransaksi(rows).Render(c.Request.Context(), c.Writer)
	})
}

func pageOv(o Overview) pages.AdminOverviewData {
	return pages.AdminOverviewData{Users: o.Users, Products: o.Products,
		PendingOrder: o.PendingOrder, DoneOrder: o.DoneOrder, OpenReports: o.OpenReports}
}
