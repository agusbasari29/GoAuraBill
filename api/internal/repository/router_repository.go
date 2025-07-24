package repository

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"gorm.io/gorm"
)

type RouterRepository interface {
	Create(router *model.Router) error
	GetAll() ([]model.Router, error)
	GetByID(id uint) (*model.Router, error)
	Update(router *model.Router) error
	Delete(id uint) error
}

type routerRepository struct {
	db *gorm.DB
}

func NewRouterRepository(db *gorm.DB) RouterRepository {
	return &routerRepository{db: db}
}

func (r *routerRepository) Create(router *model.Router) error {
	return r.db.Create(router).Error
}

func (r *routerRepository) GetAll() ([]model.Router, error) {
	var routers []model.Router
	err := r.db.Find(&routers).Error
	return routers, err
}

func (r *routerRepository) GetByID(id uint) (*model.Router, error) {
	var router model.Router
	err := r.db.First(&router, id).Error
	return &router, err
}

func (r *routerRepository) Update(router *model.Router) error {
	return r.db.Save(router).Error
}

func (r *routerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Router{}, id).Error
}