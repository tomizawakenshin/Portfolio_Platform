// 4. backend/adapters/persistence/post_repo.go
package portfolio

import (
	"backend/domain/portfolio"
	domainUser "backend/domain/user"

	"gorm.io/gorm"
)

// postRepo は domain/portfolio.Repository の具象実装です
type postRepo struct {
	db *gorm.DB
}

// NewPostRepo は GORMを使ったリポジトリ実装を生成します
func NewPostRepo(db *gorm.DB) portfolio.Repository {
	return &postRepo{db: db}
}

// Create はドメインモデルを永続化モデルにマッピングして保存します
func (r *postRepo) CreatePost(p *portfolio.Post) error {
	pm := PostModel{
		Title:       p.Title,
		Description: p.Description,
		Genres:      p.Genres,
		Skills:      p.Skills,
		UserID:      p.UserID,
	}
	for _, img := range p.Images {
		pm.Images = append(pm.Images, ImageModel{URL: img.URL})
	}
	return r.db.Create(&pm).Error
}

// FindByID は GORMから取得したモデルをドメインモデルに変換します
func (r *postRepo) GetPostByID(id uint) (*portfolio.Post, error) {
	var pm PostModel
	if err := r.db.
		Preload("User").
		Preload("Images").
		First(&pm, id).Error; err != nil {
		return nil, err
	}
	imgs := make([]portfolio.Image, len(pm.Images))
	for i, im := range pm.Images {
		imgs[i] = portfolio.Image{URL: im.URL}
	}

	du := domainUser.UserModel{
		ID:               pm.User.ID,
		FirstName:        pm.User.FirstName,
		LastName:         pm.User.LastName,
		FirstNameKana:    pm.User.FirstNameKana,
		LastNameKana:     pm.User.LastNameKana,
		Email:            pm.User.Email,
		SchoolName:       pm.User.SchoolName,
		Department:       pm.User.Department,
		Laboratory:       pm.User.Laboratory,
		GraduationYear:   pm.User.GraduationYear,
		DesiredJobTypes:  []string(pm.User.DesiredJobTypes),
		Skills:           []string(pm.User.Skills),
		SelfIntroduction: pm.User.SelfIntroduction,
		ProfileImageURL:  pm.User.ProfileImageURL,
	}
	return &portfolio.Post{
		ID:          pm.ID,
		Title:       pm.Title,
		Description: pm.Description,
		Genres:      pm.Genres,
		Skills:      pm.Skills,
		Images:      imgs,
		UserID:      pm.UserID,
		User:        du,
		CreatedAt:   pm.CreatedAt,
		UpdatedAt:   pm.UpdatedAt,
	}, nil
}

// FindByUserID はユーザーIDで絞り込み、結果をドメインモデルにマッピングします
func (r *postRepo) GetPostsByUserID(userID uint) ([]*portfolio.Post, error) {
	var pms []PostModel
	if err := r.db.Where("user_id = ?", userID).Preload("Images").Find(&pms).Error; err != nil {
		return nil, err
	}
	var posts []*portfolio.Post
	for _, pm := range pms {
		imgs := make([]portfolio.Image, len(pm.Images))
		for i, im := range pm.Images {
			imgs[i] = portfolio.Image{URL: im.URL}
		}
		posts = append(posts, &portfolio.Post{
			ID:          pm.ID,
			Title:       pm.Title,
			Description: pm.Description,
			Genres:      pm.Genres,
			Skills:      pm.Skills,
			Images:      imgs,
			UserID:      pm.UserID,
			CreatedAt:   pm.CreatedAt,
			UpdatedAt:   pm.UpdatedAt,
		})
	}
	return posts, nil
}

// GetAllPosts は全件取得し、必ず空のスライスを返すようにします
func (r *postRepo) GetAllPosts() ([]*portfolio.Post, error) {
	var pms []PostModel
	if err := r.db.
		Preload("User"). // ← ここを追加
		Preload("Images").
		Find(&pms).
		Error; err != nil {
		return nil, err
	}
	// ここで空スライスを初期化することで、後段で JSON にシリアライズするとき
	// nil ではなく [] になります
	posts := make([]*portfolio.Post, 0, len(pms))
	for i := range pms {
		posts = append(posts, toDomain(&pms[i]))
	}
	return posts, nil
}

// toDomain は PostModel → domain.Post へのマッピング関数です
func toDomain(pm *PostModel) *portfolio.Post {
	imgs := make([]portfolio.Image, len(pm.Images))
	for i, im := range pm.Images {
		imgs[i] = portfolio.Image{URL: im.URL}
	}

	user := domainUser.UserModel{
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
	}

	return &portfolio.Post{
		ID:          pm.ID,
		Title:       pm.Title,
		Description: pm.Description,
		Genres:      pm.Genres,
		Skills:      pm.Skills,
		Images:      imgs,
		UserID:      pm.UserID,
		User:        user,
		CreatedAt:   pm.CreatedAt,
		UpdatedAt:   pm.UpdatedAt,
	}
}
