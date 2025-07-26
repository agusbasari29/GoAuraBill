package model
import (
	"time"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	UserID      uint          `gorm:"not null;uniqueIndex"`
	User        User          `gorm:"foreignKey:UserID"`
	Address     string        `gorm:"type:text"`
	Phone       string        `gorm:"size:20;index"`
	IDCard      string        `gorm:"size:50;uniqueIndex"` // Nomor KTP/SIM
	ProfileID   uint          // Paket langganan
	Profile     ServiceProfile `gorm:"foreignKey:ProfileID"`
	ExpiryDate  time.Time     // Tanggal berakhir langganan
	Balance     float64       `gorm:"type:decimal(10,2);default:0"`
	Status      string        `gorm:"size:20;default:'active'"` // active, suspended, terminated
}