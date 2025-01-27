package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type IGenreRepository interface {
	GetAll() ([]models.Genre, error)
}

type GenreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(db *gorm.DB) IGenreRepository {
	return &GenreRepository{db: db}
}

func (r *GenreRepository) GetAll() ([]models.Genre, error) {
	var genre []models.Genre
	if err := r.db.Find(&genre).Error; err != nil {
		return nil, err
	}
	return genre, nil
}
