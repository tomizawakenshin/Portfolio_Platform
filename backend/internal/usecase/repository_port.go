package usecase

import "backend/internal/domain"

// RepositoryPort はポートフォリオ永続化の抽象です
type IRepositoryPort interface {
	CreatePost(p *domain.Post) error
	GetPostByID(id uint) (*domain.Post, error)
	GetPostsByUserID(userID uint) ([]*domain.Post, error)
	GetAllPosts() ([]*domain.Post, error)
}
