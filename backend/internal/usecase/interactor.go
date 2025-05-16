// internal/portfolio/usecase/interactor.go
package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"backend/dto"
	"backend/internal/domain"
)

// portfolioInteractor は PortfolioUsecase を実装します
type portfolioInteractor struct {
	repo IRepositoryPort
}

// NewInteractor はユースケース実装を返します
func NewInteractor(r IRepositoryPort) PortfolioUsecase {
	return &portfolioInteractor{repo: r}
}

// CreatePost は投稿の作成を行います。
// ──DTO＋ファイル群→画像保存→ドメインエンティティ生成→永続化
func (u *portfolioInteractor) CreatePost(input dto.CreatePostInput, files []*multipart.FileHeader, userID uint) error {
	// 1) 画像を保存して domain.Image スライスを作る
	var images []domain.Image
	for _, fh := range files {
		if fh.Size > 8*1024*1024 {
			return fmt.Errorf("file %s is too large", fh.Filename)
		}
		img, err := saveImage(fh)
		if err != nil {
			return err
		}
		images = append(images, img)
	}

	// 2) ドメインエンティティを生成 (バリデーション含む)
	post, err := domain.NewPost(
		input.Title,
		input.Description,
		input.Genres,
		input.Skills,
		images,
		userID,
	)
	if err != nil {
		return err
	}

	// 3) リポジトリに保存
	return u.repo.CreatePost(post)
}

// GetPostByID は単一の投稿を取得します
func (u *portfolioInteractor) GetPostByID(id uint) (*domain.Post, error) {
	return u.repo.GetPostByID(id)
}

// GetPostsByUserID はユーザーIDで投稿を絞り込みます
func (u *portfolioInteractor) GetPostsByUserID(userID uint) ([]*domain.Post, error) {
	return u.repo.GetPostsByUserID(userID)
}

// GetAllPosts は全投稿を取得します
func (u *portfolioInteractor) GetAllPosts() ([]*domain.Post, error) {
	return u.repo.GetAllPosts()
}

// saveImage はファイルをディスクに保存し、domain.Image を返します
func saveImage(fh *multipart.FileHeader) (domain.Image, error) {
	file, err := fh.Open()
	if err != nil {
		return domain.Image{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	uploadDir := "uploads/PortfolioImages"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return domain.Image{}, fmt.Errorf("failed to create upload dir: %v", err)
		}
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fh.Filename)
	path := fmt.Sprintf("%s/%s", uploadDir, filename)

	out, err := os.Create(path)
	if err != nil {
		return domain.Image{}, fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return domain.Image{}, fmt.Errorf("failed to save file: %v", err)
	}

	return domain.Image{URL: path}, nil
}
