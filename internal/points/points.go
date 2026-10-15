package points

// Aturan perolehan dan biaya tukar sesuai DESIGN Bab Poin & Dampak.
// Level selain "Pemula Hijau", ambang, dan kedaluwarsa masih TBD.
const (
	EarnBeli        = 10
	EarnJualLengkap = 15
	EarnTanpaRetur  = 5
	EarnRework      = 20
	EarnReferral    = 25

	CostOngkir = 100
	CostDonasi = 80
	CostPohon  = 120
)

type Action string

const (
	ActionBeli        Action = "beli"
	ActionJualLengkap Action = "jual_lengkap"
	ActionTanpaRetur  Action = "tanpa_retur"
	ActionRework      Action = "rework"
	ActionReferral    Action = "referral"
)

type Reward string

const (
	RewardOngkir Reward = "ongkir"
	RewardDonasi Reward = "donasi"
	RewardPohon  Reward = "pohon"
)

// EarnFor mengembalikan delta poin untuk aksi. ok=false bila aksi tak dikenal.
func EarnFor(a Action) (delta int, ok bool) {
	switch a {
	case ActionBeli:
		return EarnBeli, true
	case ActionJualLengkap:
		return EarnJualLengkap, true
	case ActionTanpaRetur:
		return EarnTanpaRetur, true
	case ActionRework:
		return EarnRework, true
	case ActionReferral:
		return EarnReferral, true
	default:
		return 0, false
	}
}

// RedeemCost mengembalikan biaya poin untuk reward. ok=false bila tak dikenal.
func RedeemCost(r Reward) (cost int, ok bool) {
	switch r {
	case RewardOngkir:
		return CostOngkir, true
	case RewardDonasi:
		return CostDonasi, true
	case RewardPohon:
		return CostPohon, true
	default:
		return 0, false
	}
}

// Apply menambah delta ke saldo. Saldo tidak boleh negatif (dependen dicek
// via CanRedeem sebelum redeem; earn selalu positif).
func Apply(balance, delta int) int {
	return balance + delta
}

// CanRedeem true bila saldo cukup untuk biaya.
func CanRedeem(balance, cost int) bool {
	return balance >= cost
}

// Shortage mengembalikan kekurangan poin (0 bila cukup).
func Shortage(balance, cost int) int {
	if balance >= cost {
		return 0
	}
	return cost - balance
}
