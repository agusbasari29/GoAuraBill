package model

import (
	"time"
	"gorm.io/gorm"
)

type Voucher struct {
	gorm.Model
	Code      string    `gorm:"uniqueIndex;not null"`
	Status    string    `gorm:"type:varchar(20);not null;default:'new'"` // 'new', 'used', 'expired'
	ProfileID uint      `gorm:"not null"`
	Profile   ServiceProfile `gorm:"foreignKey:ProfileID"`
	UsedBy    *uint     // User ID yang menggunakan voucher
	UsedAt    *time.Time
}