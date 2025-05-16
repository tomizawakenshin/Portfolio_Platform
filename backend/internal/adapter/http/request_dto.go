package http

import (
	"backend/dto"
	"mime/multipart"
)

// CreatePostRequest はポートフォリオ作成エンドポイントのリクエストパラメータを表します
// Gin の `ShouldBind` / `MultipartForm` でバインド可能なタグを指定します
// 複数の画像ファイルは `images` フィールドで受け取ります

type CreatePostRequest struct {
	Title       string                  `form:"title" binding:"required"`
	Description string                  `form:"description" binding:"required"`
	Genres      []string                `form:"genres" binding:"required"`
	Skills      []string                `form:"skills" binding:"required"`
	Images      []*multipart.FileHeader `form:"images" binding:"required"`
}

// ToDTO は CreatePostRequest をサービス層に渡す DTO に変換します
func (r *CreatePostRequest) ToDTO(userID uint) (out dto.CreatePostInput, files []*multipart.FileHeader) {
	out = dto.CreatePostInput{
		Title:       r.Title,
		Description: r.Description,
		Genres:      r.Genres,
		Skills:      r.Skills,
	}
	files = r.Images
	return
}
