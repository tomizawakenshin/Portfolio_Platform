// services/job_type_service.go

package services

import (
	"backend/models"
	"backend/repositories"
)

type IJobTypeService interface {
	GetAllJobTypes() ([]models.JobType, error)
}

type JobTypeService struct {
	jobTypeRepository repositories.IJobTypeRepository
}

func NewJobTypeService(jobTypeRepository repositories.IJobTypeRepository) IJobTypeService {
	return &JobTypeService{jobTypeRepository: jobTypeRepository}
}

func (s *JobTypeService) GetAllJobTypes() ([]models.JobType, error) {
	return s.jobTypeRepository.GetAll()
}
