package services

import (
	"backend/models"
	reposotories "backend/repositories"
)

type IUserService interface {
	GetUserByID(userID uint) (*models.User, error)
}

type UserService struct {
	repository reposotories.IUserRepository
}

func NewUserService(repository reposotories.IUserRepository) IUserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	return s.repository.FindByID(userID)
}
