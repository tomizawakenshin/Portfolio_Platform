package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName         string
	LastName          string
	Email             string `gorm:"not null;unique"`
	Password          string `gorm:"not null"`
	IsVerified        bool   `gorm:"default:false"` // 仮登録フラグ
	VerificationToken string `gorm:"not null"`      // 本登録用トークン
}
