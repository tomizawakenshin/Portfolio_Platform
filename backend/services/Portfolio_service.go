package services

import (
	"backend/models"
	"backend/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type IPortfolioService interface {
	CreatePost(ctx *gin.Context, userID uint) error
	GetPostByID(id uint) (*models.Post, error)
	GetPostsByUserID(userID uint) ([]models.Post, error)
	GetAllPosts() ([]models.Post, error)
}

type PortfolioService struct {
	portfolioRepository repositories.IPortfolioRepository
}

func NewPortfolioService(portfolioRepository repositories.IPortfolioRepository) IPortfolioService {
	return &PortfolioService{portfolioRepository: portfolioRepository}
}

func (s *PortfolioService) CreatePost(ctx *gin.Context, userID uint) error {
	// テキストフィールドを取得
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")

	// ジャンルを取得
	genres := ctx.PostFormArray("genres")
	fmt.Println("Current genres:", genres)
	// スキルとジャンルを取得（JSON文字列をパースする必要があるかもしれません）
	skills := ctx.PostFormArray("skills")

	// ファイルを取得
	form, _ := ctx.MultipartForm()
	files := form.File["images"]

	var images []models.Image

	for _, fileHeader := range files {
		// ファイルサイズの制限をチェック（例: 8MB）
		if fileHeader.Size > 8*1024*1024 {
			return fmt.Errorf("file %s is too large", fileHeader.Filename)
		}

		// ファイルを保存し、画像モデルを作成
		image, err := saveImage(fileHeader)
		if err != nil {
			return err
		}
		images = append(images, image)
	}

	// 投稿データを作成
	post := models.Post{
		Title:       title,
		Description: description,
		Genres:      genres,
		Skills:      skills,
		UserID:      userID,
		Images:      images,
	}

	// リポジトリを呼び出して投稿を作成
	return s.portfolioRepository.CreatePost(&post)
}

func (s *PortfolioService) GetPostByID(id uint) (*models.Post, error) {
	return s.portfolioRepository.GetPostByID(id)
}

func (s *PortfolioService) GetPostsByUserID(userID uint) ([]models.Post, error) {
	return s.portfolioRepository.GetPostsByUserID(userID)
}

// 画像を保存し、Imageモデルを返す
func saveImage(fileHeader *multipart.FileHeader) (models.Image, error) {
	// ファイルを開く
	file, err := fileHeader.Open()
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// アップロード先ディレクトリを確認または作成
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			return models.Image{}, fmt.Errorf("failed to create upload directory: %v", err)
		}
	}

	// ユニークなファイル名を生成
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)

	// 保存先のパスを設定
	savePath := fmt.Sprintf("%s/%s", uploadDir, filename)

	// ファイルを保存
	out, err := os.Create(savePath)
	if err != nil {
		return models.Image{}, fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return models.Image{}, fmt.Errorf("failed to save file: %v", err)
	}

	// 画像モデルを作成
	image := models.Image{
		// IDはDBによって自動生成されるため、IDフィールドは不要
		URL: savePath, // 必要に応じてURLに変換
	}

	return image, nil
}

func (s *PortfolioService) GetAllPosts() ([]models.Post, error) {
	return s.portfolioRepository.GetAllPosts()
}
