package service

import (
	"errors"
	"time"

	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
)

type TransactionService interface {
	CreateTransaction(txn *model.Transaction) error
	GetTransaction(id uint) (*model.Transaction, error)
	GetCustomerTransactions(customerID uint) ([]model.Transaction, error)
	ProcessPayment(txnID uint, refID string) error
	CancelTransaction(txnID uint) error
	GetPendingTransactions() ([]model.Transaction, error)
}
type transactionService struct {
	repo         repository.TransactionRepository
	customerRepo repository.CustomerRepository
}

func NewTransactionService(repo repository.TransactionRepository, customerRepo repository.CustomerRepository) TransactionService {
	return &transactionService{repo, customerRepo}
}
func (s *transactionService) CreateTransaction(txn *model.Transaction) error {
	// Validasi dasar
	if txn.CustomerID == 0 {
		return errors.New("customer ID is required")
	}
	if txn.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	// Set default status jika belum di-set
	if txn.Status == "" {
		txn.Status = model.TransactionStatusPending
	}
	return s.repo.Create(txn)
}
func (s *transactionService) GetTransaction(id uint) (*model.Transaction, error) {
	return s.repo.GetByID(id)
}
func (s *transactionService) GetCustomerTransactions(customerID uint) ([]model.Transaction, error) {
	return s.repo.GetByCustomer(customerID)
}
func (s *transactionService) ProcessPayment(txnID uint, refID string) error {
	txn, err := s.repo.GetByID(txnID)
	if err != nil {
		return err
	}
	// Validasi status
	if txn.Status != model.TransactionStatusPending {
		return errors.New("only pending transactions can be processed")
	}
	// Update status dan reference
	txn.Status = model.TransactionStatusCompleted
	txn.ReferenceID = refID
	now := time.Now()
	txn.ProcessedAt = &now
	// Simpan perubahan
	return s.repo.Update(txn)
}
func (s *transactionService) CancelTransaction(txnID uint) error {
	txn, err := s.repo.GetByID(txnID)
	if err != nil {
		return err
	}
	// Validasi status
	if txn.Status != model.TransactionStatusPending {
		return errors.New("only pending transactions can be cancelled")
	}
	txn.Status = model.TransactionStatusCancelled
	return s.repo.Update(txn)
}
func (s *transactionService) GetPendingTransactions() ([]model.Transaction, error) {
	return s.repo.GetByStatus(model.TransactionStatusPending)
}
