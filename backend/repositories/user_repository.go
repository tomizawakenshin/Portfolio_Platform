// repositories/user_repository.go

package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error // 追加
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
