package repository

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetByID(id uint) (*model.Transaction, error)
	GetByReference(ref string) (*model.Transaction, error)
	GetByCustomer(customerID uint) ([]model.Transaction, error)
	GetByType(txnType model.TransactionType) ([]model.Transaction, error)
	GetByStatus(status model.TransactionStatus) ([]model.Transaction, error)
	Update(transaction *model.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetByID(id uint) (*model.Transaction, error) {
	var txn model.Transaction
	err := r.db.Preload("User").Preload("Customer").
		First(&txn, id).Error
	return &txn, err
}

func (r *transactionRepository) GetByReference(ref string) (*model.Transaction, error) {
	var txn model.Transaction
	err := r.db.Preload("User").Preload("Customer").
		Where("reference_id = ?", ref).First(&txn).Error
	return &txn, err
}

func (r *transactionRepository) GetByCustomer(customerID uint) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.db.Preload("User").Preload("Customer").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").Find(&txns).Error
	return txns, err
}

func (r *transactionRepository) GetByType(txnType model.TransactionType) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.db.Preload("User").Preload("Customer").
		Where("type = ?", txnType).
		Order("created_at DESC").Find(&txns).Error
	return txns, err
}

func (r *transactionRepository) GetByStatus(status model.TransactionStatus) ([]model.Transaction, error) {
	var txns []model.Transaction
	err := r.db.Preload("User").Preload("Customer").
		Where("status = ?", status).
		Order("created_at DESC").Find(&txns).Error
	return txns, err
}

func (r *transactionRepository) Update(transaction *model.Transaction) error {
	return r.db.Save(transaction).Error
}
