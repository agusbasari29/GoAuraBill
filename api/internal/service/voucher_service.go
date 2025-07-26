package service

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
	"github.com/agusbasari29/GoAuraBill/internal/util"
)

type VoucherService interface {
	GenerateVouchers(quantity int, profileID uint) ([]*model.Voucher, error)
	GetAllVouchers() ([]model.Voucher, error)
	GetVoucherByID(id uint) (*model.Voucher, error)
	DeleteVoucher(id uint) error
}

type voucherService struct {
	repo repository.VoucherRepository
}

func NewVoucherService(repo repository.VoucherRepository) VoucherService {
	return &voucherService{repo: repo}
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
func (s *voucherService) DeleteVoucher(id uint) error {
	return s.repo.Delete(id)
}
