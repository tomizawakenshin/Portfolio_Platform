// models/skill.go

package models

import (
	"gorm.io/gorm"
)

type Skill struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}
