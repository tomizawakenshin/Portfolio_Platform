// services/user_service.go

package services

import (
	domainUser "backend/domain/user"
	"backend/dto"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

type IUserService interface {
	GetUserByID(userID uint) (*domainUser.UserModel, error)
	UpdateMinimumUserInfo(userID uint, input dto.MinimumUserInfoInput, files []*multipart.FileHeader) (*domainUser.UserModel, error)
}

type UserService struct {
	repository domainUser.IUserRepository
}

func NewUserService(repository domainUser.IUserRepository) IUserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetUserByID(userID uint) (*domainUser.UserModel, error) {
	return s.repository.FindByID(userID)
}

func (s *UserService) UpdateMinimumUserInfo(userID uint, input dto.MinimumUserInfoInput, files []*multipart.FileHeader) (*domainUser.UserModel, error) {
	// DBからユーザーを取得
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 各フィールドが nil でなければ上書き
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.FirstNameKana != nil {
		user.FirstNameKana = *input.FirstNameKana
	}
	if input.LastNameKana != nil {
		user.LastNameKana = *input.LastNameKana
	}
	if input.SchoolName != nil {
		user.SchoolName = *input.SchoolName
	}
	if input.Department != nil {
		user.Department = *input.Department
	}
	if input.Laboratory != nil {
		user.Laboratory = *input.Laboratory
	}
	if input.GraduationYear != nil {
		user.GraduationYear = *input.GraduationYear
	}
	if input.DesiredJobTypes != nil {
		user.DesiredJobTypes = *input.DesiredJobTypes
	}
	if input.Skills != nil {
		user.Skills = *input.Skills
	}
	if input.SelfIntroduction != nil {
		user.SelfIntroduction = *input.SelfIntroduction
	}

	if len(files) > 0 {
		fileHeader := files[0] // 一枚だけの場合
		if fileHeader.Size > 8*1024*1024 {
			return nil, fmt.Errorf("file %s is too large", fileHeader.Filename)
		}

		// 実際の保存ロジックをこのサービス内(またはprivate関数)で呼ぶ
		savedPath, err := saveUserImage(fileHeader)
		if err != nil {
			return nil, err
		}
		user.ProfileImageURL = savedPath
	}

	// データベースに保存
	if err := s.repository.UpdateUser(user); err != nil {
		return nil, err
	}

	// 成功時、更新後の user を返す
	return user, nil
}

func saveUserImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 保存先フォルダ
	uploadDir := "uploads/UserImages"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %v", err)
		}
	}

	// ユニークファイル名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	savePath := fmt.Sprintf("%s/%s", uploadDir, filename)

	out, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return savePath, nil
}
