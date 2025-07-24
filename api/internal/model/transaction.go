package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID      uint    `gorm:"not null"`
	Amount      float64 `gorm:"type:decimal(10,2);not null"`
	Type        string  `gorm:"type:varchar(20);not null"` // 'topup', 'payment', 'refund'
	Description string
	ReferenceID string  // ID dari payment gateway atau referensi manual
}