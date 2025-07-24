package service

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5" // Pastikan Anda menginstal ini: go get github.com/golang-jwt/jwt/v5
	"github.com/agusbasari29/GoAuraBill/internal/model" // Ganti dengan path modul Anda
	"github.com/agusbasari29/GoAuraBill/internal/repository" // Ganti dengan path modul Anda
	"golang.org/x/crypto/bcrypt"
)
type AuthService interface {
	RegisterUser(user *model.User) error
	LoginUser(username, password string) (string, error) // Mengembalikan token JWT
}
type authService struct {
	authRepo  repository.AuthRepository
	jwtSecret string
}
func NewAuthService(authRepo repository.AuthRepository, jwtSecret string) AuthService {
	return &authService{authRepo: authRepo, jwtSecret: jwtSecret}
}
func (s *authService) RegisterUser(user *model.User) error {
	// Password hashing dilakukan di BeforeSave hook model User
	return s.authRepo.CreateUser(user)
}
func (s *authService) LoginUser(username, password string) (string, error) {
	user, err := s.authRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials") // Jangan berikan detail error spesifik
	}
	// Bandingkan password yang diberikan dengan hash yang tersimpan
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	// Buat token JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return signedToken, nil
}