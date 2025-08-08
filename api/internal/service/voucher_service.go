package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
	"github.com/agusbasari29/GoAuraBill/internal/util"
	"gorm.io/gorm"
)

type VoucherService interface {
	GenerateVouchers(quantity int, profileID uint) ([]*model.Voucher, error)
	GetAllVouchers() ([]model.Voucher, error)
	GetVoucherByID(id uint) (*model.Voucher, error)
	DeleteVoucher(id uint) error
	RedeemVoucher(code string, customerID uint) error
}

type voucherService struct {
	repo         repository.VoucherRepository
	customerRepo repository.CustomerRepository
	transRepo    repository.TransactionRepository
	db           *gorm.DB
}

func NewVoucherService(
	repo repository.VoucherRepository,
	customerRepo repository.CustomerRepository,
	transRepo repository.TransactionRepository,
	db *gorm.DB,
) VoucherService {
	return &voucherService{
		repo:         repo,
		customerRepo: customerRepo,
		transRepo:    transRepo,
		db:           db,
	}
}

func (s *voucherService) GenerateVouchers(quantity int, profileID uint) ([]*model.Voucher, error) {
	var vouchers []*model.Voucher
	for i := 0; i < quantity; i++ {
		voucher := &model.Voucher{
			Code:      util.GenerateRandomCode(8), // Membuat kode 8 karakter
			Status:    "new",
			ProfileID: profileID,
		}
		vouchers = append(vouchers, voucher)
	}
	err := s.repo.CreateBatch(vouchers)
	if err != nil {
		return nil, err
	}
	return vouchers, nil
}

func (s *voucherService) GetAllVouchers() ([]model.Voucher, error) {
	return s.repo.GetAll()
}

func (s *voucherService) GetVoucherByID(id uint) (*model.Voucher, error) {
	return s.repo.GetByID(id)
}

func (s *voucherService) RedeemVoucher(code string, customerID uint) error {
	// Memulai transaksi database
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Cari voucher berdasarkan kode di dalam transaksi
		voucher, err := s.repo.GetByCode(code)
		if err != nil {
			return errors.New("voucher tidak ditemukan")
		}

		// 2. Validasi status voucher
		if voucher.Status != "new" {
			return errors.New("voucher sudah digunakan atau tidak valid")
		}

		// 3. Cari data pelanggan
		customer, err := s.customerRepo.GetCustomerByID(customerID)
		if err != nil {
			return errors.New("pelanggan tidak ditemukan")
		}

		// 4. Terapkan benefit ke pelanggan
		// Mengganti paket dan memperpanjang masa aktif
		customer.ProfileID = voucher.ProfileID
		// Tambahkan masa aktif sesuai paket dari voucher
		customer.ExpiryDate = time.Now().AddDate(0, 0, voucher.Profile.ValidityDays)
		if err := s.customerRepo.UpdateCustomer(customer); err != nil {
			return fmt.Errorf("gagal memperbarui profil pelanggan: %w", err)
		}

		// 5. Perbarui status voucher
		now := time.Now()
		voucher.Status = "used"
		voucher.UsedBy = &customer.UserID
		voucher.UsedAt = &now
		if err := s.repo.Update(voucher); err != nil {
			return fmt.Errorf("gagal memperbarui status voucher: %w", err)
		}

		// 6. Catat transaksi aktivasi voucher
		transaction := &model.Transaction{
			UserID:      customer.UserID,
			CustomerID:  customer.ID,
			Amount:      voucher.Profile.Price, // Ambil harga dari profil voucher
			Type:        model.TransactionTypeVoucher,
			Status:      model.TransactionStatusCompleted,
			Description: fmt.Sprintf("Aktivasi paket '%s' menggunakan voucher %s", voucher.Profile.Name, voucher.Code),
			ReferenceID: fmt.Sprintf("VCR-%d", voucher.ID),
			ProcessedAt: &now,
		}
		if err := s.transRepo.Create(transaction); err != nil {
			return fmt.Errorf("gagal mencatat transaksi: %w", err)
		}

		// Jika semua berhasil, commit transaksi
		return nil
	})
}

func (s *voucherService) DeleteVoucher(id uint) error {
	return s.repo.Delete(id)
}
