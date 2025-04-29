package portfolio

import (
	"time"

	userInfra "backend/infrastructure/user"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// PostModel は GORM タグ付きの永続化用モデルです
// ドメインモデルとは分離し、ORM依存を閉じ込めます
type PostModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title       string              `gorm:"not null"`
	Description string              `gorm:"type:text"`
	Genres      pq.StringArray      `gorm:"type:text[]"`
	Skills      pq.StringArray      `gorm:"type:text[]"`
	Images      []ImageModel        `gorm:"foreignKey:PostID"`
	UserID      uint                `gorm:"not null;index"`
	User        userInfra.UserModel `gorm:"foreignKey:UserID;references:ID"`
}

// ImageModel は永続化層の画像モデルです
type ImageModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	URL    string `gorm:"not null"`
	PostID uint   `gorm:"not null;index"`
}
