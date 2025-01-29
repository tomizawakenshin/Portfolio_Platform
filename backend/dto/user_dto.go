package dto

type MinimumUserInfoInput struct {
	FirstName        *string   `json:"firstName"`
	LastName         *string   `json:"lastName"`
	FirstNameKana    *string   `json:"firstNameKana"`
	LastNameKana     *string   `json:"lastNameKana"`
	SchoolName       *string   `json:"schoolName"`
	Department       *string   `json:"department"`
	Laboratory       *string   `json:"laboratory"`
	GraduationYear   *string   `json:"graduationYear"`
	DesiredJobTypes  *[]string `json:"desiredJobTypes"`
	Skills           *[]string `json:"skills"`
	SelfIntroduction *string   `json:"selfIntroduction"`
}
