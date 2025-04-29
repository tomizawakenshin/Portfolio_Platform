package services

import (
	domainPortfolio "backend/domain/portfolio"
	"backend/dto"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

type IPortfolioService interface {
	CreatePost(input dto.CreatePostInput, files []*multipart.FileHeader, userID uint) error
	GetPostByID(id uint) (*domainPortfolio.Post, error)
	GetPostsByUserID(userID uint) ([]*domainPortfolio.Post, error)
	GetAllPosts() ([]*domainPortfolio.Post, error)
}

type PortfolioService struct {
	// portfolioRepository repositories.IPortfolioRepository
	portfolioRepository domainPortfolio.Repository
}

func NewPortfolioService(portfolioRepository domainPortfolio.Repository) IPortfolioService {
	return &PortfolioService{portfolioRepository: portfolioRepository}
}

func (s *PortfolioService) CreatePost(input dto.CreatePostInput,
	files []*multipart.FileHeader,
	userID uint) error {

	// 1) 画像を保存
	var images []domainPortfolio.Image
	for _, fileHeader := range files {
		if fileHeader.Size > 8*1024*1024 {
			return fmt.Errorf("file %s is too large", fileHeader.Filename)
		}
		image, err := saveImage(fileHeader)
		if err != nil {
			return err
		}
		images = append(images, domainPortfolio.Image{
			URL: image.URL,
		})
	}

	post, err := domainPortfolio.NewPost(
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
	return s.portfolioRepository.CreatePost(post)
}
func (s *PortfolioService) GetPostByID(id uint) (*domainPortfolio.Post, error) {
	return s.portfolioRepository.GetPostByID(id)
}

func (s *PortfolioService) GetPostsByUserID(userID uint) ([]*domainPortfolio.Post, error) {
	return s.portfolioRepository.GetPostsByUserID(userID)
}

// 画像を保存し、Imageモデルを返す
func saveImage(fileHeader *multipart.FileHeader) (domainPortfolio.Image, error) {
	// ファイルを開く
	file, err := fileHeader.Open()
	if err != nil {
		return domainPortfolio.Image{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// アップロード先ディレクトリを確認または作成
	uploadDir := "uploads/PortfolioImages"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			return domainPortfolio.Image{}, fmt.Errorf("failed to create upload directory: %v", err)
		}
	}

	// ユニークなファイル名を生成
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)

	// 保存先のパスを設定
	savePath := fmt.Sprintf("%s/%s", uploadDir, filename)

	// ファイルを保存
	out, err := os.Create(savePath)
	if err != nil {
		return domainPortfolio.Image{}, fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return domainPortfolio.Image{}, fmt.Errorf("failed to save file: %v", err)
	}

	// 画像モデルを作成
	image := domainPortfolio.Image{
		// IDはDBによって自動生成されるため、IDフィールドは不要
		URL: savePath, // 必要に応じてURLに変換
	}

	return image, nil
}

func (s *PortfolioService) GetAllPosts() ([]*domainPortfolio.Post, error) {
	return s.portfolioRepository.GetAllPosts()
}
