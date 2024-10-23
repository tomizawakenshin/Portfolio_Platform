// repositories/job_type_repository.go

package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type IJobTypeRepository interface {
	GetAll() ([]models.JobType, error)
}

type JobTypeRepository struct {
	db *gorm.DB
}

func NewJobTypeRepository(db *gorm.DB) IJobTypeRepository {
	return &JobTypeRepository{db: db}
}

func (r *JobTypeRepository) GetAll() ([]models.JobType, error) {
	var jobTypes []models.JobType
	if err := r.db.Find(&jobTypes).Error; err != nil {
		return nil, err
	}
	return jobTypes, nil
}
