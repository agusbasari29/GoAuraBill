package repository

import (
	"github.com/agusbasari29/GoAuraBill/internal/model" // Ganti dengan path modul Anda
	"gorm.io/gorm"
)
type AuthRepository interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
}
type authRepository struct {
	db *gorm.DB
}
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
func (r *authRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}