// dto/post_dto.go

package dto

type CreatePostInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Genres      []string `json:"genres" binding:"required"`
	Skills      []string `json:"skills"`
}
