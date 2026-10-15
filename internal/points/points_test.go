package points

import (
	"testing"
)

func TestEarnTable(t *testing.T) {
	cases := []struct {
		action Action
		want   int
	}{
		{ActionBeli, 10},
		{ActionJualLengkap, 15},
		{ActionTanpaRetur, 5},
		{ActionRework, 20},
		{ActionReferral, 25},
	}
	for _, c := range cases {
		got, ok := EarnFor(c.action)
		if !ok || got != c.want {
			t.Errorf("EarnFor(%q) = %d,%v; mau %d,true", c.action, got, ok, c.want)
		}
	}
	if _, ok := EarnFor("ngarang"); ok {
		t.Error("aksi tak dikenal harus ok=false")
	}
}

func TestRedeemTable(t *testing.T) {
	cases := []struct {
		reward Reward
		want   int
	}{
		{RewardOngkir, 100},
		{RewardDonasi, 80},
		{RewardPohon, 120},
	}
	for _, c := range cases {
		got, ok := RedeemCost(c.reward)
		if !ok || got != c.want {
			t.Errorf("RedeemCost(%q) = %d,%v; mau %d,true", c.reward, got, ok, c.want)
		}
	}
	if _, ok := RedeemCost("ngarang"); ok {
		t.Error("reward tak dikenal harus ok=false")
	}
}

func TestRedeemFlow(t *testing.T) {
	balance := Apply(0, EarnBeli) // +10
	balance = Apply(balance, EarnJualLengkap)
	if balance != 25 {
		t.Fatalf("saldo = %d, mau 25", balance)
	}
	if CanRedeem(balance, CostDonasi) {
		t.Error("saldo 25 tidak boleh tukar donasi 80")
	}
	if Shortage(balance, CostDonasi) != 55 {
		t.Errorf("kekurangan mau 55")
	}
	balance = Apply(balance, 100)
	if !CanRedeem(balance, CostDonasi) {
		t.Error("saldo 125 harus bisa tukar donasi 80")
	}
	balance = Apply(balance, -CostDonasi)
	if balance != 45 {
		t.Errorf("saldo akhir = %d, mau 45", balance)
	}
}
