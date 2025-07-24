package model

import "gorm.io/gorm"

type ServiceProfile struct {
	gorm.Model
	Name         string  `gorm:"uniqueIndex;not null"`
	DownloadRate uint    `gorm:"not null"` // dalam kbps
	UploadRate   uint    `gorm:"not null"` // dalam kbps
	Price        float64 `gorm:"type:decimal(10,2);not null;default:0.00"`
	ValidityDays int     `gorm:"not null;default:30"`
}