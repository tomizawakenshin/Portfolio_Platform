package http

import (
	"backend/internal/domain"
	"strconv"
	"time"
)

type UserResponse struct {
	ID               uint     `json:"id"`
	FirstName        string   `json:"firstName"`
	LastName         string   `json:"lastName"`
	FirstNameKana    string   `json:"firstNameKana"`
	LastNameKana     string   `json:"lastNameKana"`
	Email            string   `json:"email"`
	SchoolName       string   `json:"schoolName"`
	Department       string   `json:"department"`
	Laboratory       string   `json:"laboratory"`
	GraduationYear   int      `json:"graduationYear"`
	DesiredJobTypes  []string `json:"desiredJobTypes"`
	Skills           []string `json:"skills"`
	SelfIntroduction string   `json:"selfIntroduction"`
	ProfileImageURL  string   `json:"profileImageUrl"`
}

// PostResponse: クライアントに返すポートフォリオ情報
type PostResponse struct {
	ID          uint            `json:"id"`
	CreatedAt   string          `json:"createdAt"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Genres      []string        `json:"genres"`
	Skills      []string        `json:"skills"`
	Images      []ImageResponse `json:"images"`
	UserID      uint            `json:"userId"`
	User        UserResponse    `json:"user"`
}

type ImageResponse struct {
	ID     uint   `json:"id"`
	URL    string `json:"url"`
	PostID uint   `json:"postId"`
}

// domain.Post → PostResponse へのマッピング関数
func NewPostResponse(p *domain.Post) *PostResponse {
	images := make([]ImageResponse, len(p.Images))
	for i, img := range p.Images {
		images[i] = ImageResponse{
			ID:     img.ID, // domain.ImageにID,PostID追加
			URL:    img.URL,
			PostID: img.PostID,
		}
	}
	// domain.Post.User → UserResponse へマッピング
	user := UserResponse{
		ID:               p.User.ID,
		FirstName:        p.User.FirstName,
		LastName:         p.User.LastName,
		FirstNameKana:    p.User.FirstNameKana,
		LastNameKana:     p.User.LastNameKana,
		Email:            p.User.Email,
		SchoolName:       p.User.SchoolName,
		Department:       p.User.Department,
		Laboratory:       p.User.Laboratory,
		GraduationYear:   graduationYearStringToInt(p.User.GraduationYear),
		DesiredJobTypes:  p.User.DesiredJobTypes,
		Skills:           p.User.Skills,
		SelfIntroduction: p.User.SelfIntroduction,
		ProfileImageURL:  p.User.ProfileImageURL,
	}

	return &PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Genres:      p.Genres,
		Skills:      p.Skills,
		Images:      images,
		UserID:      p.UserID,
		User:        user,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
	}
}

// 複数形（配列用）
func NewPostsResponse(posts []*domain.Post) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, p := range posts {
		responses[i] = NewPostResponse(p)
	}
	return responses
}

// 補助関数：string→int（失敗時は0返す）
func graduationYearStringToInt(str string) int {
	year, err := strconv.Atoi(str)
	if err != nil {
		return 0 // エラー時は 0 を返す
	}
	return year
}
