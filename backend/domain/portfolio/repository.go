// 2. backend/domain/portfolio/repository.go
package portfolio

// Repository は投稿エンティティの永続化を抽象化したインターフェースです
// サービス層はこのインターフェースだけを依存先として扱います
type Repository interface {
	CreatePost(post *Post) error
	GetPostByID(id uint) (*Post, error)
	GetPostsByUserID(userID uint) ([]*Post, error)
	GetAllPosts() ([]*Post, error)
}
