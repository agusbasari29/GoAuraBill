package service

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
)

type CustomerService interface {
	CreateCustomer(customer *model.User) error
	GetAllCustomers() ([]model.User, error)
	GetCustomerByID(id uint) (*model.User, error)
	UpdateCustomer(customer *model.User) error
	DeleteCustomer(id uint) error
}

type customerService struct {
	repo repository.CustomerRepository
	// Di masa depan, kita akan menambahkan service Mikrotik di sini
	// mikrotikService MikrotikService
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(customer *model.User) error {
	// Logika bisnis sebelum membuat customer:
	// 1. Validasi apakah ProfileID valid (ada di DB).
	// 2. Panggil service Mikrotik untuk membuat user di router.
	// Untuk saat ini, kita langsung memanggil repository.
	return s.repo.Create(customer)
}

func (s *customerService) GetAllCustomers() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *customerService) GetCustomerByID(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *customerService) UpdateCustomer(customer *model.User) error {
	// Logika bisnis sebelum update:
	// 1. Validasi data.
	// 2. Panggil service Mikrotik untuk update user di router.
	return s.repo.Update(customer)
}

func (s *customerService) DeleteCustomer(id uint) error {
	// Logika bisnis sebelum delete:
	// 1. Panggil service Mikrotik untuk menghapus user dari router.
	return s.repo.Delete(id)
}