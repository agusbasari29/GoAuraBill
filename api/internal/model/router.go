package model

import "gorm.io/gorm"

type Router struct {
	gorm.Model
	Name     string `gorm:"not null"`
	IPAddress string `gorm:"uniqueIndex;not null"`
	Port     string `gorm:"not null;default:'8728'"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"` // Sebaiknya dienkripsi saat disimpan
	IsActive bool   `gorm:"default:true"`
}
