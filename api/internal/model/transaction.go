package model
import (
	"time"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeTopUp    TransactionType = "topup"
	TransactionTypePayment  TransactionType = "payment"
	TransactionTypeRefund    TransactionType = "refund"
	TransactionTypeAdjust   TransactionType = "adjustment"
	TransactionTypeVoucher  TransactionType = "voucher"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

type Transaction struct {
	gorm.Model
	UserID      uint            `gorm:"not null"`
	User        User            `gorm:"foreignKey:UserID"`
	CustomerID  uint            `gorm:"not null"`
	Customer    Customer        `gorm:"foreignKey:CustomerID"`
	Amount      float64         `gorm:"type:decimal(10,2);not null"`
	Type        TransactionType `gorm:"type:varchar(20);not null"`
	Status      TransactionStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	Description string
	ReferenceID string          `gorm:"index"` // ID dari payment gateway atau referensi eksternal
	Metadata    JSON            `gorm:"type:json"` // Data tambahan dalam format JSON
	ProcessedAt *time.Time      // Waktu transaksi diproses
}

type JSON map[string]interface{}
// Hook AfterCreate untuk transaksi otomatis
func (t *Transaction) AfterCreate(tx *gorm.DB) (err error) {
	if t.Type == TransactionTypeTopUp && t.Status == TransactionStatusCompleted {
		// Update saldo customer jika transaksi topup berhasil
		err = tx.Model(&Customer{}).Where("id = ?", t.CustomerID).
			Update("balance", gorm.Expr("balance + ?", t.Amount)).Error
	}
	return
}