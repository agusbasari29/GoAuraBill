package repository

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"gorm.io/gorm"
)

type VoucherRepository interface {
	CreateBatch(vouchers []*model.Voucher) error
	GetAll() ([]model.Voucher, error)
	GetByID(id uint) (*model.Voucher, error)
	GetByCode(code string) (*model.Voucher, error)
	Update(voucher *model.Voucher) error
	Delete(id uint) error
}

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db: db}
}
func (r *voucherRepository) CreateBatch(vouchers []*model.Voucher) error {
	// Menggunakan CreateInBatches untuk efisiensi saat membuat banyak data
	return r.db.CreateInBatches(vouchers, 100).Error
}
func (r *voucherRepository) GetAll() ([]model.Voucher, error) {
	var vouchers []model.Voucher
	// Mengambil data voucher dengan preload relasi Profile
	err := r.db.Preload("Profile").Find(&vouchers).Error
	return vouchers, err
}
func (r *voucherRepository) GetByID(id uint) (*model.Voucher, error) {
	var voucher model.Voucher
	err := r.db.Preload("Profile").First(&voucher, id).Error
	return &voucher, err
}
func (r *voucherRepository) Delete(id uint) error {
	return r.db.Delete(&model.Voucher{}, id).Error
}

func (r *voucherRepository) GetByCode(code string) (*model.Voucher, error) {
	var voucher model.Voucher
	// Preload Profile untuk mendapatkan detail paket dari voucher
	err := r.db.Preload("Profile").Where("code = ?", code).First(&voucher).Error
	return &voucher, err
}

func (r *voucherRepository) Update(voucher *model.Voucher) error {
	return r.db.Save(voucher).Error
}