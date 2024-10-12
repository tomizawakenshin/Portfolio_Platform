package reposotories

import (
	"backend/models"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(user models.User) error
	FindUser(email string) (*models.User, error)
	FindUserByToken(token string) (*models.User, error)
	UpdateUser(user *models.User) error
	SoftDeleteUnverifiedUsersBefore(cutoffTime time.Time) error
	PermanentlyDeleteUsersBefore(cutoffTime time.Time) error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user models.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AuthRepository) FindUser(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByToken(token string) (*models.User, error) {
	var user models.User
	result := r.db.Where("verification_token = ?", token).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
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
