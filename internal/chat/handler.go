package chat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/models"
	"mewali/web/templates/components"
	"mewali/web/templates/pages"
)

// RegisterRoutes memasang mulai thread, isi thread, kirim, dan partial polling.
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	login := auth.RequireLogin(gdb, secret)

	r.POST("/chat/mulai", login, func(c *gin.Context) {
		u := currentMember(c)
		pid, _ := strconv.ParseUint(c.PostForm("product_id"), 10, 64)
		var p models.Product
		if err := gdb.First(&p, uint(pid)).Error; err != nil {
			c.String(http.StatusNotFound, "produk tidak ditemukan")
			return
		}
		if p.SellerID == u.ID {
			c.String(http.StatusBadRequest, "tidak bisa chat dengan diri sendiri")
			return
		}
		th, err := GetOrCreateThread(gdb, p.ID, u.ID, p.SellerID, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "thread gagal")
			return
		}
		c.Redirect(http.StatusFound, "/chat/"+strconv.FormatUint(uint64(th.ID), 10))
	})

	r.GET("/chat/:id", login, func(c *gin.Context) {
		th, u, ok := memberThread(c, gdb)
		if !ok {
			return
		}
		msgs, _ := ListMessages(gdb, th.ID, 100)
		var peer models.User
		peerID := th.BuyerID
		if u.ID == th.BuyerID {
			peerID = th.SellerID
		}
		gdb.First(&peer, peerID)
		d := pages.ChatData{Thread: th, Messages: msgs, PeerName: peer.Name, Me: u.ID}
		c.Status(http.StatusOK)
		_ = pages.Chat(d).Render(c.Request.Context(), c.Writer)
	})

	r.POST("/chat/:id", login, func(c *gin.Context) {
		th, u, ok := memberThread(c, gdb)
		if !ok {
			return
		}
		if err := PostMessage(gdb, th.ID, u.ID, c.PostForm("body")); err != nil {
			c.String(http.StatusBadRequest, "gagal: %v", err)
			return
		}
		c.Redirect(http.StatusFound, "/chat/"+strconv.FormatUint(uint64(th.ID), 10))
	})

	r.GET("/partial/chat/:id", login, func(c *gin.Context) {
		th, u, ok := memberThread(c, gdb)
		if !ok {
			return
		}
		msgs, _ := ListMessages(gdb, th.ID, 100)
		c.Status(http.StatusOK)
		_ = components.ChatList(msgs, u.ID).Render(c.Request.Context(), c.Writer)
	})
}

func currentMember(c *gin.Context) models.User {
	u, _ := c.Get("user")
	return u.(models.User)
}

func memberThread(c *gin.Context, gdb *gorm.DB) (models.ChatThread, models.User, bool) {
	var th models.ChatThread
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := gdb.First(&th, uint(id)).Error; err != nil {
		c.String(http.StatusNotFound, "chat tidak ditemukan")
		return th, models.User{}, false
	}
	u := currentMember(c)
	if u.ID != th.BuyerID && u.ID != th.SellerID {
		c.String(http.StatusForbidden, "bukan chatmu")
		return th, u, false
	}
	return th, u, true
}
