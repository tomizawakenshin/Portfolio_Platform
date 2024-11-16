// dto/post_dto.go

package dto

type CreatePostInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	GitHubLink  string   `json:"githubLink"`
	ProductLink string   `json:"productLink"`
	Skills      []string `json:"skills"`
}
