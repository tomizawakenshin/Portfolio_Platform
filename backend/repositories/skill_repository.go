// repositories/skill_repository.go

package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type ISkillRepository interface {
	GetAll() ([]models.Skill, error)
}

type SkillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) ISkillRepository {
	return &SkillRepository{db: db}
}

func (r *SkillRepository) GetAll() ([]models.Skill, error) {
	var skills []models.Skill
	if err := r.db.Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}
