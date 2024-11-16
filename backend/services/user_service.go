// services/user_service.go

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
	user.FirstNameKana = input.FirstNameKana
	user.LastNameKana = input.LastNameKana
	user.SchoolName = input.SchoolName
	user.Department = input.Department
	user.Laboratory = input.Laboratory
	user.GraduationYear = input.GraduationYear
	user.DesiredJobTypes = input.DesiredJobTypes
	user.Skills = input.Skills

	// ユーザーを保存
	if err := s.repository.UpdateUser(user); err != nil {
		return err
	}

	return nil
}
