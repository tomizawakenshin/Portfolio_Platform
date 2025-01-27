// services/genre_service.go

package services

import (
	"backend/models"
	"backend/repositories"
)

type IGenreService interface {
	GetAllGenre() ([]models.Genre, error)
}

type GenreService struct {
	genreRepository repositories.IGenreRepository
}

func NewGenreService(genreRepository repositories.IGenreRepository) IGenreService {
	return &GenreService{genreRepository: genreRepository}
}

func (s *GenreService) GetAllGenre() ([]models.Genre, error) {
	return s.genreRepository.GetAll()
}
