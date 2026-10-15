package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/models"
	"mewali/web/templates/pages"
)

const sessionTTL = 30 * 24 * time.Hour

// RegisterRoutes memasang /login dan /logout.
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	r.GET("/login", func(c *gin.Context) {
		if _, ok := currentUser(c, gdb, secret); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.Status(http.StatusOK)
		_ = pages.Login("").Render(c.Request.Context(), c.Writer)
	})
	r.POST("/login", func(c *gin.Context) {
		u, err := Login(gdb, c.PostForm("email"), c.PostForm("password"))
		if err != nil {
			c.Status(http.StatusUnauthorized)
			_ = pages.Login("Email atau sandi salah. Coba lagi ya.").Render(c.Request.Context(), c.Writer)
			return
		}
		c.SetCookie(CookieName, Issue(secret, u.ID, sessionTTL, time.Now()),
			int(sessionTTL.Seconds()), "/", "", false, true)
		c.Redirect(http.StatusFound, "/")
	})
	r.POST("/logout", func(c *gin.Context) {
		c.SetCookie(CookieName, "", -1, "/", "", false, true)
		c.Redirect(http.StatusFound, "/login")
	})
}

// RequireLogin menolak anonim ke /login. User tersimpan di c.Get("user").
func RequireLogin(gdb *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := currentUser(c, gdb, secret)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("user", u)
		c.Next()
	}
}

func currentUser(c *gin.Context, gdb *gorm.DB, secret string) (models.User, bool) {
	uid, ok := CurrentUserID(c, secret)
	if !ok {
		return models.User{}, false
	}
	u, err := FindByID(gdb, uid)
	if err != nil {
		return models.User{}, false
	}
	return u, true
}

// CurrentUserID memverifikasi cookie sesi dan mengembalikan ID user.
// Diekspor agar handler modul lain (katalog, dsb.) bisa personalisasi.
func CurrentUserID(c *gin.Context, secret string) (uint, bool) {
	cookie, err := c.Cookie(CookieName)
	if err != nil {
		return 0, false
	}
	return Verify(secret, cookie, time.Now())
}
