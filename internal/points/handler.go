package points

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mewali/internal/auth"
	"mewali/internal/models"
	"mewali/web/templates/pages"
)

// RegisterRoutes memasang /poin dan /redeem.
func RegisterRoutes(r *gin.Engine, gdb *gorm.DB, secret string) {
	login := auth.RequireLogin(gdb, secret)

	r.GET("/poin", login, func(c *gin.Context) {
		u := currentViewer(c)
		renderPoin(c, gdb, u, "")
	})

	r.POST("/redeem", login, func(c *gin.Context) {
		u := currentViewer(c)
		if err := Redeem(gdb, u.ID, Reward(c.PostForm("reward"))); err != nil {
			renderPoin(c, gdb, u, redeemMsg(err))
			return
		}
		c.Redirect(http.StatusFound, "/poin")
	})
}

func currentViewer(c *gin.Context) models.User {
	u, _ := c.Get("user")
	return u.(models.User)
}

func renderPoin(c *gin.Context, gdb *gorm.DB, u models.User, errMsg string) {
	var fresh models.User
	gdb.First(&fresh, u.ID)
	balance := Balance(gdb, u.ID)
	rewards := []struct {
		r    Reward
		name string
	}{
		{RewardOngkir, "Gratis Ongkir"},
		{RewardDonasi, "Donasi 1 pakaian layak"},
		{RewardPohon, "Tanam 1 pohon"},
	}
	d := pages.PoinData{Balance: balance, Level: fresh.Level, Error: errMsg}
	for _, rw := range rewards {
		cost, _ := RedeemCost(rw.r)
		d.Rewards = append(d.Rewards, pages.RewardView{
			Reward: string(rw.r), Name: rw.name, Cost: cost,
			CanAfford: CanRedeem(balance, cost), Shortage: Shortage(balance, cost),
		})
	}
	c.Status(http.StatusOK)
	_ = pages.Poin(d).Render(c.Request.Context(), c.Writer)
}

func redeemMsg(err error) string {
	if err == ErrInsufficient {
		return "Poin belum cukup untuk reward ini."
	}
	return "Gagal: " + err.Error()
}
