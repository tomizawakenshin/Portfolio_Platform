package user

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// UserModel は GORM タグ付きの永続化用ユーザーモデルです
// ドメインモデルとは分離し、ORM依存を閉じ込めます
// このファイルは infrastructure/user/model.go に配置します

type UserModel struct {
	// ID        uint `gorm:"primaryKey"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	gorm.Model

	FirstName             string  `gorm:"not null"`
	LastName              string  `gorm:"not null"`
	FirstNameKana         string  `gorm:"not null"`
	LastNameKana          string  `gorm:"not null"`
	Email                 string  `gorm:"not null;unique"`
	Password              *string `gorm:"size:255"`
	IsVerified            bool    `gorm:"default:false"`
	VerificationToken     *string `gorm:"size:255"`
	VerificationExpiresAt time.Time
	PasswordResetToken    string `gorm:"size:255"`
	PasswordResetExpires  time.Time

	SchoolName     string `gorm:"size:255"`
	Department     string `gorm:"size:255"`
	Laboratory     string `gorm:"size:255"`
	GraduationYear string `gorm:"size:4"`

	DesiredJobTypes pq.StringArray `gorm:"type:text[]"`
	Skills          pq.StringArray `gorm:"type:text[]"`

	SelfIntroduction string `gorm:"type:text"`
	ProfileImageURL  string `gorm:"size:512"`
}
