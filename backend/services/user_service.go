package services

import (
	"backend/dto"
	"backend/models"
	reposotories "backend/repositories"
)

type IUserService interface {
	GetUserByID(userID uint) (*models.User, error)
	UpdateMinimumUserInfo(userID uint, input dto.MinimumUserInfoInput) error
}

type UserService struct {
	repository     reposotories.IUserRepository
	authRepository reposotories.IAuthRepository
}

func NewUserService(repository reposotories.IUserRepository, authRepository reposotories.IAuthRepository) IUserService {
	return &UserService{
		repository:     repository,
		authRepository: authRepository,
	}
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	return s.repository.FindByID(userID)
}

func (s *UserService) UpdateMinimumUserInfo(userID uint, input dto.MinimumUserInfoInput) error {
	// ユーザーを取得
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}

	// ユーザー情報を更新
	user.FirstName = &input.FirstName
	user.LastName = &input.LastName

	// ユーザーを保存
	if err := s.authRepository.UpdateUser(user); err != nil {
		return err
	}

	return nil
}
