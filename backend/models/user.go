// models/user.go

package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName             string
	LastName              string
	FirstNameKana         string
	LastNameKana          string
	Email                 string `gorm:"not null;unique"`
	Password              *string
	IsVerified            bool `gorm:"default:false"`
	VerificationToken     *string
	VerificationExpiresAt time.Time
	PasswordResetToken    string
	PasswordResetExpires  time.Time
	SchoolName            string
	Department            string
	Laboratory            string
	GraduationYear        string
	DesiredJobTypes       pq.StringArray `gorm:"type:text[]"`
	Skills                pq.StringArray `gorm:"type:text[]"`
}
