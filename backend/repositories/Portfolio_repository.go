package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type IPortfolioRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id uint) (*models.Post, error)
	GetPostsByUserID(userID uint) ([]models.Post, error)
	GetAllPosts() ([]models.Post, error)
}

type PortfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) IPortfolioRepository {
	return &PortfolioRepository{db: db}
}

func (r *PortfolioRepository) CreatePost(post *models.Post) error {
	// トランザクションを使用して投稿と画像を保存
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		for _, image := range post.Images {
			image.PostID = post.ID
			image.ID = 0
			if err := tx.Create(&image).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PortfolioRepository) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Images").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PortfolioRepository) GetPostsByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Where("user_id = ?", userID).Preload("Images").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PortfolioRepository) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Preload("Images").Preload("User").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
