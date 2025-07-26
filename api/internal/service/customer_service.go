package service
import (
	"errors"
	"time"
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
)

type CustomerService interface {
	CreateCustomer(customer *model.Customer) error
	GetCustomerByID(id uint) (*model.Customer, error)
	GetCustomerByUserID(userID uint) (*model.Customer, error)
	GetAllCustomers() ([]model.Customer, error)
	GetCustomersByStatus(status string) ([]model.Customer, error)
	UpdateCustomer(customer *model.Customer) error
	SuspendCustomer(id uint) error
	ActivateCustomer(id uint) error
	TerminateCustomer(id uint) error
	DeleteCustomer(id uint) error
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo}
}

func (s *customerService) CreateCustomer(customer *model.Customer) error {
	// Validasi data customer
	if customer.UserID == 0 {
		return errors.New("user ID is required")
	}
	if customer.ProfileID == 0 {
		return errors.New("service profile is required")
	}
	
	// Set default expiry date (30 hari dari sekarang)
	if customer.ExpiryDate.IsZero() {
		customer.ExpiryDate = time.Now().AddDate(0, 0, 30)
	}
	
	return s.repo.CreateCustomer(customer)
}

func (s *customerService) GetCustomerByID(id uint) (*model.Customer, error) {
	return s.repo.GetCustomerByID(id)
}

func (s *customerService) GetCustomerByUserID(userID uint) (*model.Customer, error) {
	return s.repo.GetCustomerByUserID(userID)
}

func (s *customerService) GetAllCustomers() ([]model.Customer, error) {
	return s.repo.GetAllCustomers()
}

func (s *customerService) GetCustomersByStatus(status string) ([]model.Customer, error) {
	validStatus := map[string]bool{
		"active": true, "suspended": true, "terminated": true,
	}
	if !validStatus[status] {
		return nil, errors.New("invalid status")
	}
	return s.repo.GetCustomersByStatus(status)
}

func (s *customerService) UpdateCustomer(customer *model.Customer) error {
	// Validasi sebelum update
	if customer.ID == 0 {
		return errors.New("customer ID is required")
	}
	return s.repo.UpdateCustomer(customer)
}

func (s *customerService) SuspendCustomer(id uint) error {
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return err
	}
	customer.Status = "suspended"
	return s.repo.UpdateCustomer(customer)
}

func (s *customerService) ActivateCustomer(id uint) error {
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return err
	}
	customer.Status = "active"
	return s.repo.UpdateCustomer(customer)
}

func (s *customerService) TerminateCustomer(id uint) error {
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return err
	}
	customer.Status = "terminated"
	return s.repo.UpdateCustomer(customer)
}

func (s *customerService) DeleteCustomer(id uint) error {
	return s.repo.DeleteCustomer(id)
}