package service

import (
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
)

type ServiceProfileService interface {
	CreateProfile(profile *model.ServiceProfile) error
	GetAllProfiles() ([]model.ServiceProfile, error)
	GetProfileByID(id uint) (*model.ServiceProfile, error)
	UpdateProfile(profile *model.ServiceProfile) error
	DeleteProfile(id uint) error
}

type serviceProfileService struct {
	repo repository.ServiceProfileRepository
}

func NewServiceProfileService(repo repository.ServiceProfileRepository) ServiceProfileService {
	return &serviceProfileService{repo: repo}
}

func (s *serviceProfileService) CreateProfile(profile *model.ServiceProfile) error {
	return s.repo.Create(profile)
}

func (s *serviceProfileService) GetAllProfiles() ([]model.ServiceProfile, error) {
	return s.repo.GetAll()
}

func (s *serviceProfileService) GetProfileByID(id uint) (*model.ServiceProfile, error) {
	return s.repo.GetByID(id)
}

func (s *serviceProfileService) UpdateProfile(profile *model.ServiceProfile) error {
	return s.repo.Update(profile)
}

func (s *serviceProfileService) DeleteProfile(id uint) error {
	return s.repo.Delete(id)
}