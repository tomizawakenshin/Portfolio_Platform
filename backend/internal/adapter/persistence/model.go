// internal/portfolio/adapter/persistence/model.go
package persistence

import (
	"time"

	userInfra "backend/infrastructure/user"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// PostModel は GORM 用の永続化モデルです。
// ドメイン層の Post エンティティを DB スキーマにマッピングします。
// Frameworks & Drivers 層に対応し、ORMタグによってテーブル定義を行います。
// この構造体はドメインロジックに依存せず、DBアクセス用のフィールドのみを持ちます。

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

// ImageModel は投稿に紐づく画像情報の永続化モデルです。
// PostModel のリレーションとして定義し、foreignKey で紐付けます。

type ImageModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	URL    string `gorm:"not null"`
	PostID uint   `gorm:"not null;index"`
}
