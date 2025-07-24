package model

import (
	"github.com/agusbasari29/GoAuraBill/internal/util" // Ganti dengan path modul Anda
	"gorm.io/gorm"
)

// Variabel global untuk menyimpan kunci enkripsi
var EncryptionKey string

type Router struct {
	gorm.Model
	Name      string `gorm:"not null"`
	IPAddress string `gorm:"uniqueIndex;not null"`
	Port      string `gorm:"not null;default:'8728'"`
	Username  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	IsActive  bool   `gorm:"default:true"`
}

// BeforeSave mengenkripsi password sebelum disimpan ke DB
func (r *Router) BeforeSave(tx *gorm.DB) (err error) {
	if r.Password != "" {
		encryptedPass, err := util.Encrypt(r.Password, EncryptionKey)
		if err != nil {
			return err
		}
		r.Password = encryptedPass
	}
	return
}

// AfterFind mendekripsi password setelah diambil dari DB
func (r *Router) AfterFind(tx *gorm.DB) (err error) {
	if r.Password != "" {
		decryptedPass, err := util.Decrypt(r.Password, EncryptionKey)
		if err != nil {
			// Jangan gagalkan seluruh proses jika dekripsi gagal, mungkin password sudah terdekripsi
			// atau dalam format lama. Cukup log error atau abaikan.
			return nil
		}
		r.Password = decryptedPass
	}
	return
}