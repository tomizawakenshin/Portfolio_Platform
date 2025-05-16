package usecase

import "backend/internal/domain"

// PortfolioOutputPort はユースケース実行後の出力を受け取るインターフェースです
type PortfolioOutputPort interface {
	// 単一の Post を返す場合
	PresentPost(post *domain.Post) error

	// 複数の Post を返す場合
	PresentPosts(posts []*domain.Post) error
}
