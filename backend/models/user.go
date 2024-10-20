package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName             string
	LastName              string
	Email                 string `gorm:"not null;unique"`
	Password              *string
	IsVerified            bool    `gorm:"default:false"` // 仮登録フラグ
	VerificationToken     *string // 本登録用トークン
	VerificationExpiresAt time.Time
	PasswordResetToken    string
	PasswordResetExpires  time.Time
}
