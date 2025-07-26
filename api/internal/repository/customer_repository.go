package repository

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(customer *model.Customer) error
	GetCustomerByID(id uint) (*model.Customer, error)
	GetCustomerByUserID(userID uint) (*model.Customer, error)
	GetAllCustomers() ([]model.Customer, error)
	GetCustomersByStatus(status string) ([]model.Customer, error)
	UpdateCustomer(customer *model.Customer) error
	DeleteCustomer(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) CreateCustomer(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) GetCustomerByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Preload("User").Preload("Profile").First(&customer, id).Error
	return &customer, err
}

func (r *customerRepository) GetCustomerByUserID(userID uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Preload("User").Preload("Profile").Where("user_id = ?", userID).First(&customer).Error
	return &customer, err
}

func (r *customerRepository) GetAllCustomers() ([]model.Customer, error) {
	var customers []model.Customer
	err := r.db.Preload("User").Preload("Profile").Find(&customers).Error
	return customers, err
}

func (r *customerRepository) GetCustomersByStatus(status string) ([]model.Customer, error) {
	var customers []model.Customer
	err := r.db.Preload("User").Preload("Profile").Where("status = ?", status).Find(&customers).Error
	return customers, err
}

func (r *customerRepository) UpdateCustomer(customer *model.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) DeleteCustomer(id uint) error {
	return r.db.Delete(&model.Customer{}, id).Error
}