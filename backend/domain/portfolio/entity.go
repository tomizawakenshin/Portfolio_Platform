// 1. backend/domain/portfolio/entity.go
package portfolio

import (
	domainUser "backend/domain/user"
	"fmt"
	"time"
)

// Image は投稿に紐づく画像情報（ドメインモデル）
type Image struct {
	URL string
}

// Post は作品投稿を表すドメインエンティティ
// ビジネスルール（必須チェックなど）を内包します
type Post struct {
	ID          uint
	Title       string
	Description string
	Genres      []string
	Skills      []string
	Images      []Image
	UserID      uint
	User        domainUser.UserModel
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPost は Post を生成するファクトリメソッドです
// タイトル必須、ジャンル1つ以上、画像1枚以上のチェックを行います
func NewPost(
	title, description string,
	genres, skills []string,
	images []Image,
	userID uint,
) (*Post, error) {
	if title == "" {
		return nil, fmt.Errorf("タイトルは必須です")
	}
	if len(genres) == 0 {
		return nil, fmt.Errorf("ジャンルは1つ以上選択してください")
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("画像は少なくとも1枚必要です")
	}
	now := time.Now()
	return &Post{
		Title:       title,
		Description: description,
		Genres:      genres,
		Skills:      skills,
		Images:      images,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
