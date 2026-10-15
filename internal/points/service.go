package points

import (
	"errors"

	"gorm.io/gorm"

	"mewali/internal/models"
)

var ErrInsufficient = errors.New("poin tidak cukup")

// Balance menjumlahkan ledger user (sumber kebenaran saldo ada di backend).
func Balance(gdb *gorm.DB, userID uint) int {
	var total int64
	gdb.Model(&models.PointLedger{}).Where("user_id = ?", userID).Select("COALESCE(SUM(delta),0)").Scan(&total)
	return int(total)
}

// Award mencatat perolehan poin untuk aksi (append-only).
func Award(gdb *gorm.DB, userID uint, orderID *uint, a Action) error {
	delta, ok := EarnFor(a)
	if !ok {
		return errors.New("aksi poin tak dikenal")
	}
	return gdb.Create(&models.PointLedger{UserID: userID, OrderID: orderID,
		Delta: delta, Reason: string(a), BalanceAfter: Balance(gdb, userID) + delta}).Error
}

// Redeem menukar reward: cek saldo, catat redemption + ledger negatif, atomik.
func Redeem(gdb *gorm.DB, userID uint, r Reward) error {
	cost, ok := RedeemCost(r)
	if !ok {
		return errors.New("reward tak dikenal")
	}
	return gdb.Transaction(func(tx *gorm.DB) error {
		if Balance(tx, userID) < cost {
			return ErrInsufficient
		}
		if err := tx.Create(&models.Redemption{UserID: userID, Reward: string(r), Cost: cost, Status: "done"}).Error; err != nil {
			return err
		}
		return tx.Create(&models.PointLedger{UserID: userID, Delta: -cost,
			Reason: "redeem_" + string(r), BalanceAfter: Balance(tx, userID) - cost}).Error
	})
}
