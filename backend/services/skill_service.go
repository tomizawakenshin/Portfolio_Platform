// services/skill_service.go

package services

import (
	"backend/models"
	"backend/repositories"
)

type ISkillService interface {
	GetAllSkills() ([]models.Skill, error)
}

type SkillService struct {
	skillRepository repositories.ISkillRepository
}

func NewSkillService(skillRepository repositories.ISkillRepository) ISkillService {
	return &SkillService{skillRepository: skillRepository}
}

func (s *SkillService) GetAllSkills() ([]models.Skill, error) {
	return s.skillRepository.GetAll()
}
