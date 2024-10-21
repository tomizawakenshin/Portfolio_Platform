package services

import (
	"backend/dto"
	"backend/models"
	"backend/repositories"
)

type IUserService interface {
	GetUserByID(userID uint) (*models.User, error)
	UpdateMinimumUserInfo(userID uint, input dto.MinimumUserInfoInput) error
}

type UserService struct {
	repository     repositories.IUserRepository
	authRepository repositories.IAuthRepository
}

func NewUserService(repository repositories.IUserRepository, authRepository repositories.IAuthRepository) IUserService {
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
	user.FirstName = input.FirstName
	user.LastName = input.LastName

	// ユーザーを保存
	if err := s.authRepository.UpdateUser(user); err != nil {
		return err
	}

	return nil
}
