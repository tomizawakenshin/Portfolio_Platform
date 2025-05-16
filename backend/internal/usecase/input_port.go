// internal/portfolio/usecase/input_port.go
package usecase

import (
	"mime/multipart"

	"backend/dto"
	"backend/internal/domain"
)

// PortfolioUsecase はポートフォリオ機能のユースケース入力ポート
type PortfolioUsecase interface {
	CreatePost(input dto.CreatePostInput, files []*multipart.FileHeader, userID uint) error
	GetPostByID(id uint) (*domain.Post, error)
	GetPostsByUserID(userID uint) ([]*domain.Post, error)
	GetAllPosts() ([]*domain.Post, error)
}
