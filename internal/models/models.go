package models

import (
	"time"

	"gorm.io/gorm"
)

// Status produk, metode bayar, mode order, dan enum lain disimpan sebagai
// string agar fleksibel (level poin, kebijakan baru) tanpa migrasi berat.

const (
	ProductDraft      = "draft"
	ProductInReview   = "in_review"
	ProductActive     = "active"
	ProductSold       = "sold"
	ProductTakenDown  = "taken_down"
	PayCash           = "cash"
	PayTransferManual = "transfer_manual"
	ModeP2P           = "p2p"    // Chat + Ketemuan, selesai manual
	ModeMeetup        = "meetup" // mock online: 2 kode + poin
	OrderPending      = "pending"
	OrderConfirmed    = "confirmed"
	OrderDone         = "done"
	OrderCancelled    = "cancelled"
	OrderExpired      = "expired"
	CodeBuyer         = "buyer"
	CodeSeller        = "seller"
	RewardOngkir      = "ongkir"
	RewardDonasi      = "donasi"
	RewardPohon       = "pohon"
	ReportOpen        = "open"
	ReportInReview    = "in_review"
	ReportResolved    = "resolved"
	ReportRejected    = "rejected"
)

type User struct {
	gorm.Model
	Name          string `gorm:"size:120;not null"`
	Email         string `gorm:"size:160;uniqueIndex;not null"`
	PasswordHash  string `gorm:"size:255"`
	AvatarURL     string `gorm:"size:500"`
	Handle        string `gorm:"size:60;uniqueIndex"`
	Lat           float64
	Lng           float64
	LocationLabel string `gorm:"size:200"`
	Level         string `gorm:"size:60;default:'Pemula Hijau'"`
	EmailVerified bool
	IsAdmin       bool `gorm:"default:false"`
}

// Preferensi ukuran tubuh (cm). Nol = belum diisi.
type SizePreference struct {
	gorm.Model
	UserID   uint `gorm:"uniqueIndex;not null"`
	Dada     float64
	Pinggang float64
	Pinggul  float64
	Tinggi   float64
	Torso    float64
}

type Product struct {
	gorm.Model
	SellerID      uint   `gorm:"index;not null"`
	Name          string `gorm:"size:160;not null"`
	Category      string `gorm:"size:60;index"`
	Price         int    `gorm:"not null"` // rupiah
	Description   string `gorm:"type:text"`
	Lat           float64
	Lng           float64
	LocationLabel string `gorm:"size:200"`
	Status        string `gorm:"size:20;index;default:'draft'"`
}

// Slot foto wajib: depan, belakang, kerah, lengan, tag_samping, tag_atas.
type ProductPhoto struct {
	gorm.Model
	ProductID uint   `gorm:"index;not null"`
	Slot      string `gorm:"size:20;not null"`
	URL       string `gorm:"size:500;not null"`
	Sort      int
}

// Penanda cacat: koordinat dinormalisasi 0-1000 agar tahan resize.
type Defect struct {
	gorm.Model
	ProductID uint   `gorm:"index;not null"`
	PhotoID   uint   `gorm:"index"`
	X         int    `gorm:"not null"`
	Y         int    `gorm:"not null"`
	Label     string `gorm:"size:300;not null"`
	Sort      int
}

type ConditionScore struct {
	gorm.Model
	ProductID uint      `gorm:"uniqueIndex;not null"`
	Score     int       `gorm:"not null"`
	Method    string    `gorm:"size:60"` // manual|aturan|... (TBD, tanpa klaim AI)
	Note      string    `gorm:"size:500"`
	CheckedAt time.Time `gorm:"not null"`
}

// Ukuran barang (cm). Nol = tidak diukur. 1:1 dengan Product.
type ProductSize struct {
	gorm.Model
	ProductID uint `gorm:"uniqueIndex;not null"`
	Dada      float64
	Pinggang  float64
	Pinggul   float64
	Tinggi    float64
	Torso     float64
}

type CartItem struct {
	gorm.Model
	UserID    uint `gorm:"index;not null;uniqueIndex:idx_cart_user_product"`
	ProductID uint `gorm:"not null;uniqueIndex:idx_cart_user_product"`
}

type Order struct {
	gorm.Model
	BuyerID   uint   `gorm:"index;not null"`
	SellerID  uint   `gorm:"index;not null"`
	ProductID uint   `gorm:"index;not null"`
	Total     int    `gorm:"not null"`
	PayMethod string `gorm:"size:20;not null"`
	Mode      string `gorm:"size:20;not null"`
	Status    string `gorm:"size:20;index;not null"`
}

type ConfirmCode struct {
	gorm.Model
	OrderID   uint      `gorm:"index;not null;uniqueIndex:idx_code_order_role"`
	Role      string    `gorm:"size:10;not null;uniqueIndex:idx_code_order_role"`
	Code      string    `gorm:"size:5;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	UsedAt    *time.Time
	Attempts  int `gorm:"default:0"`
}

// Ledger append-only: jangan update/delete baris.
type PointLedger struct {
	gorm.Model
	UserID       uint   `gorm:"index;not null"`
	OrderID      *uint  `gorm:"index"`
	Delta        int    `gorm:"not null"`
	Reason       string `gorm:"size:120;not null"`
	BalanceAfter int    `gorm:"not null"`
}

type Redemption struct {
	gorm.Model
	UserID uint   `gorm:"index;not null"`
	Reward string `gorm:"size:20;not null"`
	Cost   int    `gorm:"not null"`
	Status string `gorm:"size:20;not null"`
}

type Review struct {
	gorm.Model
	OrderID    uint   `gorm:"uniqueIndex;not null"`
	ReviewerID uint   `gorm:"index;not null"`
	SellerID   uint   `gorm:"index;not null"`
	Stars      int    `gorm:"not null"`
	Body       string `gorm:"type:text"`
}

type Report struct {
	gorm.Model
	ReporterID   uint   `gorm:"index;not null"`
	TargetUserID *uint  `gorm:"index"`
	ProductID    *uint  `gorm:"index"`
	Reason       string `gorm:"size:120;not null"`
	Detail       string `gorm:"type:text"`
	Status       string `gorm:"size:20;index;default:'open'"`
}

type ChatThread struct {
	gorm.Model
	OrderID   *uint `gorm:"index"`
	ProductID uint  `gorm:"index;not null"`
	BuyerID   uint  `gorm:"index;not null"`
	SellerID  uint  `gorm:"index;not null"`
}

type ChatMessage struct {
	gorm.Model
	ThreadID uint   `gorm:"index;not null"`
	SenderID uint   `gorm:"index;not null"`
	Body     string `gorm:"type:text;not null"`
}

type DonationPoint struct {
	gorm.Model
	Name         string `gorm:"size:160;not null"`
	Lat          float64
	Lng          float64
	AddressLabel string `gorm:"size:300"`
	Active       bool   `gorm:"default:true"`
}

type Banner struct {
	gorm.Model
	Title    string `gorm:"size:160;not null"`
	Subtitle string `gorm:"size:300"`
	CTALabel string `gorm:"size:60"`
	CTAHref  string `gorm:"size:300"`
	Sort     int
	Active   bool `gorm:"default:true"`
}

type AdminAudit struct {
	gorm.Model
	AdminID    uint   `gorm:"index;not null"`
	Action     string `gorm:"size:120;not null"`
	TargetType string `gorm:"size:60"`
	TargetID   uint
	Note       string `gorm:"size:500"`
}

// Kategori awal sesuai DESIGN Bab chip kategori. Di models (bukan catalog)
// agar template bisa pakai tanpa siklus impor.
var Categories = []string{"Kaos", "Kemeja", "Jeans", "Gaun", "Jaket"}

// All mendaftarkan semua model untuk AutoMigrate (dev) dan referensi
// migrasi SQL prod. Urutan: tabel tanpa FK dulu tidak wajib (GORM bebas).
func All() []any {
	return []any{
		&User{}, &SizePreference{},
		&Product{}, &ProductPhoto{}, &Defect{}, &ConditionScore{}, &ProductSize{},
		&CartItem{}, &Order{}, &ConfirmCode{},
		&PointLedger{}, &Redemption{},
		&Review{}, &Report{},
		&ChatThread{}, &ChatMessage{},
		&DonationPoint{}, &Banner{}, &AdminAudit{},
	}
}
