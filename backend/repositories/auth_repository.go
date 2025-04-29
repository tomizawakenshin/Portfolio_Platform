package repositories

import (
	"backend/models"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(user models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByVerificationToken(token string) (*models.User, error)
	UpdateUser(user *models.User) error
	SoftDeleteUnverifiedUsersBefore(cutoffTime time.Time) error
	PermanentlyDeleteUsersBefore(cutoffTime time.Time) error
	FindUserByPasswordResetToken(token string) (*models.User, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{db: db}
}

const (
	fieldEmail              = "email"
	fieldVerificationToken  = "verification_token"
	fieldPasswordResetToken = "password_reset_token"
)

func (r *AuthRepository) CreateUser(user models.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	return r.FindUserByField(fieldEmail, email)
}

func (r *AuthRepository) FindUserByVerificationToken(token string) (*models.User, error) {
	return r.FindUserByField(fieldVerificationToken, token)
}

func (r *AuthRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *AuthRepository) FindUserByField(fieldName string, value interface{}) (*models.User, error) {
	var user models.User
	result := r.db.Where(fmt.Sprintf("%s = ?", fieldName), value).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// ソフトデリートを行うメソッドを追加
func (r *AuthRepository) SoftDeleteUnverifiedUsersBefore(cutoffTime time.Time) error {
	result := r.db.Where("is_verified = ? AND created_at < ?", false, cutoffTime).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Soft-deleted %d unverified users.", result.RowsAffected)
	return nil
}

// ハードデリートを行うメソッドを追加
func (r *AuthRepository) PermanentlyDeleteUsersBefore(cutoffTime time.Time) error {
	result := r.db.Unscoped().Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffTime).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Permanently deleted %d users.", result.RowsAffected)
	return nil
}

func (r *AuthRepository) FindUserByPasswordResetToken(token string) (*models.User, error) {
	return r.FindUserByField(fieldPasswordResetToken, token)
}
