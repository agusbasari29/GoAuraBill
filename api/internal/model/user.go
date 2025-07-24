package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName    string        `gorm:"not null"`
	Username    string        `gorm:"uniqueIndex;not null"`
	Email       string        `gorm:"uniqueIndex;not null"`
	Password    string        `gorm:"not null"`
	Role        string        `gorm:"type:varchar(20);not null;default:'customer'"` // 'admin' or 'customer'
	IsActive    bool          `gorm:"default:true"`
	Balance     float64       `gorm:"type:decimal(10,2);not null;default:0.00"`
	ProfileID   *uint         // Foreign key untuk ServiceProfile (nullable)
	Profile     ServiceProfile `gorm:"foreignKey:ProfileID"`
	Transactions []Transaction `gorm:"foreignKey:UserID"`
}

// Hook BeforeSave untuk melakukan hash pada password sebelum disimpan
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return
}