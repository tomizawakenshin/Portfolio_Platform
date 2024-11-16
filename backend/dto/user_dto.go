// dto/minimum_user_info_input.go

package dto

type MinimumUserInfoInput struct {
	FirstName       string   `json:"firstName" binding:"required"`
	LastName        string   `json:"lastName" binding:"required"`
	FirstNameKana   string   `json:"firstNameKana" binding:"required"`
	LastNameKana    string   `json:"lastNameKana" binding:"required"`
	SchoolName      string   `json:"schoolName" binding:"required"`
	Department      string   `json:"department" binding:"required"`
	Laboratory      string   `json:"laboratory" binding:"required"`
	GraduationYear  string   `json:"graduationYear" binding:"required"`
	DesiredJobTypes []string `json:"desiredJobTypes" binding:"required"`
	Skills          []string `json:"skills"` // スキルは任意
}
