package repository

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"gorm.io/gorm"
)

type ServiceProfileRepository interface {
	Create(profile *model.ServiceProfile) error
	GetAll() ([]model.ServiceProfile, error)
	GetByID(id uint) (*model.ServiceProfile, error)
	Update(profile *model.ServiceProfile) error
	Delete(id uint) error
}

type serviceProfileRepository struct {
	db *gorm.DB
}

func NewServiceProfileRepository(db *gorm.DB) ServiceProfileRepository {
	return &serviceProfileRepository{db: db}
}

func (r *serviceProfileRepository) Create(profile *model.ServiceProfile) error {
	return r.db.Create(profile).Error
}

func (r *serviceProfileRepository) GetAll() ([]model.ServiceProfile, error) {
	var profiles []model.ServiceProfile
	err := r.db.Find(&profiles).Error
	return profiles, err
}

func (r *serviceProfileRepository) GetByID(id uint) (*model.ServiceProfile, error) {
	var profile model.ServiceProfile
	err := r.db.First(&profile, id).Error
	return &profile, err
}

func (r *serviceProfileRepository) Update(profile *model.ServiceProfile) error {
	return r.db.Save(profile).Error
}

func (r *serviceProfileRepository) Delete(id uint) error {
	return r.db.Delete(&model.ServiceProfile{}, id).Error
}