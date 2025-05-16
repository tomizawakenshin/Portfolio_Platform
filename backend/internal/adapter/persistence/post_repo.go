// internal/portfolio/adapter/persistence/post_repo.go
package persistence

import (
	domainUser "backend/domain/user"
	"backend/internal/domain"
	"backend/internal/usecase"

	"gorm.io/gorm"
)

// postRepo は domain.RepositoryPort インターフェースの GORM 実装です。
// Interface Adapters 層（Gateways）に対応し、ドメインリポジトリの具体実装を提供します。
// Usecase 層はこの実装を知らず、インターフェースを通して永続化を行います。

type postRepo struct {
	db *gorm.DB
}

// NewPostRepo は GORM DB を使ったリポジトリ実装を返します。
// 依存性注入（DI）によりコンストラクタで提供します。

func NewPostRepo(db *gorm.DB) usecase.IRepositoryPort {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(p *domain.Post) error {
	// ドメインモデル ↔ 永続化モデルのマッピング
	pm := PostModel{
		Title:       p.Title,
		Description: p.Description,
		Genres:      p.Genres,
		Skills:      p.Skills,
		UserID:      p.UserID,
	}
	if err := r.db.Create(&pm).Error; err != nil {
		return err
	}
	// 保存後の ID やタイムスタンプは不要であれば無視
	return nil
}

func (r *postRepo) GetPostByID(id uint) (*domain.Post, error) {
	var pm PostModel
	// 画像やユーザー情報を Preload する場合はここで設定可能
	if err := r.db.Preload("Images").First(&pm, id).Error; err != nil {
		return nil, err
	}
	// ドメインモデルに変換
	post := &domain.Post{
		ID:          pm.ID,
		Title:       pm.Title,
		Description: pm.Description,
		Genres:      pm.Genres,
		Skills:      pm.Skills,
		UserID:      pm.UserID,
		CreatedAt:   pm.CreatedAt,
		UpdatedAt:   pm.UpdatedAt,
	}
	// 画像処理: ImageModel → domain.Image
	for _, im := range pm.Images {
		post.Images = append(post.Images, domain.Image{URL: im.URL})
	}
	return post, nil
}

func (r *postRepo) GetPostsByUserID(userID uint) ([]*domain.Post, error) {
	var pms []PostModel
	if err := r.db.Where("user_id = ?", userID).Preload("Images").Find(&pms).Error; err != nil {
		return nil, err
	}
	var posts []*domain.Post
	for _, pm := range pms {
		post := &domain.Post{
			ID:          pm.ID,
			Title:       pm.Title,
			Description: pm.Description,
			Genres:      pm.Genres,
			Skills:      pm.Skills,
			UserID:      pm.UserID,
			CreatedAt:   pm.CreatedAt,
			UpdatedAt:   pm.UpdatedAt,
		}
		for _, im := range pm.Images {
			post.Images = append(post.Images, domain.Image{URL: im.URL})
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postRepo) GetAllPosts() ([]*domain.Post, error) {
	var pms []PostModel
	if err := r.db.
		Preload("User").
		Preload("Images").
		Find(&pms).Error; err != nil {
		return nil, err
	}
	posts := make([]*domain.Post, len(pms))
	for i := range pms {
		posts[i] = toDomain(&pms[i])
	}

	return posts, nil
}

// toDomain は PostModel → domain.Post へのマッピングを共通化します
func toDomain(pm *PostModel) *domain.Post {
	// ドメインエンティティを構築
	post := &domain.Post{
		ID:          pm.ID,
		Title:       pm.Title,
		Description: pm.Description,
		Genres:      pm.Genres,
		Skills:      pm.Skills,
		UserID:      pm.UserID,
		CreatedAt:   pm.CreatedAt,
		UpdatedAt:   pm.UpdatedAt,
		User: domainUser.UserModel{
			ID:               pm.User.ID,
			FirstName:        pm.User.FirstName,
			LastName:         pm.User.LastName,
			FirstNameKana:    pm.User.FirstNameKana,
			LastNameKana:     pm.User.LastNameKana,
			ProfileImageURL:  pm.User.ProfileImageURL,
			SchoolName:       pm.User.SchoolName,
			Department:       pm.User.Department,
			Laboratory:       pm.User.Laboratory,
			GraduationYear:   pm.User.GraduationYear,
			DesiredJobTypes:  []string(pm.User.DesiredJobTypes),
			Skills:           []string(pm.User.Skills),
			SelfIntroduction: pm.User.SelfIntroduction,
			CreatedAt:        pm.User.CreatedAt,
			UpdatedAt:        pm.User.UpdatedAt,
			DeletedAt:        pm.User.DeletedAt.Time,
		},
	}

	// Images をドメインモデルに変換
	for _, im := range pm.Images {
		post.Images = append(post.Images, domain.Image{URL: im.URL})
	}
	return post
}
