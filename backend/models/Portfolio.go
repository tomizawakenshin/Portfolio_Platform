// models/post.go

package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	GitHubLink  string
	ProductLink string
	Skills      pq.StringArray `gorm:"type:text[]"`
	Images      []Image        `gorm:"foreignKey:PostID"`
	UserID      uint           `gorm:"not null"`
	User        User           `gorm:"foreignKey:UserID"`
}

type Image struct {
	gorm.Model
	URL    string `gorm:"not null"`
	PostID uint   `gorm:"not null"`
}
