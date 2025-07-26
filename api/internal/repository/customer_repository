package repository

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *model.User) error
	GetAll() ([]model.User, error)
	GetByID(id uint) (*model.User, error)
	Update(customer *model.User) error
	Delete(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// Create memastikan role pengguna adalah 'customer'
func (r *customerRepository) Create(customer *model.User) error {
	customer.Role = "customer" // Paksa role menjadi customer
	return r.db.Create(customer).Error
}

// GetAll hanya mengambil pengguna dengan role 'customer'
func (r *customerRepository) GetAll() ([]model.User, error) {
	var customers []model.User
	err := r.db.Where("role = ?", "customer").Find(&customers).Error
	return customers, err
}

// GetByID hanya mengambil pengguna dengan role 'customer'
func (r *customerRepository) GetByID(id uint) (*model.User, error) {
	var customer model.User
	err := r.db.Where("role = ? AND id = ?", "customer", id).First(&customer).Error
	return &customer, err
}

func (r *customerRepository) Update(customer *model.User) error {
	customer.Role = "customer" // Pastikan role tidak berubah saat update
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Where("role = ? AND id = ?", "customer", id).Delete(&model.User{}).Error
}