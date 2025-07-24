package service

import (
	// "context"
	// "crypto/tls"
	// "errors"
	"fmt"
	"github.com/go-routeros/routeros" // go get github.com/go-mikrotik/mikrotik
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
)
type RouterService interface {
	CreateRouter(router *model.Router) error
	GetAllRouters() ([]model.Router, error)
	GetRouterByID(id uint) (*model.Router, error)
	UpdateRouter(router *model.Router) error
	DeleteRouter(id uint) error
}

type routerService struct {
	repo repository.RouterRepository
}

func NewRouterService(repo repository.RouterRepository) RouterService {
	return &routerService{repo: repo}
}

// testConnection mencoba menghubungkan ke Mikrotik menggunakan pustaka baru
func (s *routerService) testConnection(router *model.Router) error {
	address := fmt.Sprintf("%s:%s", router.IPAddress, router.Port)

	// Konfigurasi koneksi
	// Anda bisa menambahkan tls.Config jika menggunakan koneksi terenkripsi
	conn, err := routeros.Dial(address, router.Username, router.Password)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke mikrotik: %w", err)
	}
	defer conn.Close()

	// Kirim perintah sederhana untuk memverifikasi koneksi
	_, err = conn.Run("/system/resource/print")
	if err != nil {
		return fmt.Errorf("gagal menjalankan perintah di mikrotik: %w", err)
	}

	return nil
}

func (s *routerService) CreateRouter(router *model.Router) error {
	// Dekripsi password tidak diperlukan di sini karena AfterFind hook
	// hanya berjalan saat mengambil data, bukan saat membuat.
	// Kita perlu password plaintext untuk tes koneksi.
	if err := s.testConnection(router); err != nil {
		return fmt.Errorf("tes koneksi gagal: %w", err)
	}
	// Enkripsi akan terjadi secara otomatis melalui GORM hook `BeforeSave`
	return s.repo.Create(router)
}

func (s *routerService) GetAllRouters() ([]model.Router, error) {
	return s.repo.GetAll()
}

func (s *routerService) GetRouterByID(id uint) (*model.Router, error) {
	return s.repo.GetByID(id)
}

func (s *routerService) UpdateRouter(router *model.Router) error {
	// Password akan otomatis terdekripsi oleh hook `AfterFind` saat data diambil
	// dan terenkripsi kembali oleh `BeforeSave` saat disimpan.
	if err := s.testConnection(router); err != nil {
		return fmt.Errorf("tes koneksi gagal: %w", err)
	}
	return s.repo.Update(router)
}

func (s *routerService) DeleteRouter(id uint) error {
	return s.repo.Delete(id)
}