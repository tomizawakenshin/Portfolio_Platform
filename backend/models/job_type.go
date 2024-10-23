// models/job_type.go

package models

import (
	"gorm.io/gorm"
)

type JobType struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}
