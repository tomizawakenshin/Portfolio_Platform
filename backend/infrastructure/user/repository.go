package user

import (
	"time"

	domainUser "backend/domain/user"

	"gorm.io/gorm"
)

// UserRepository は domain.User.Repository の具象実装です
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepo は GORM を使った User リポジトリを生成します
func NewUserRepository(db *gorm.DB) domainUser.IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u *domainUser.UserModel) error {
	pm := toPersistence(u)
	return r.db.Create(&pm).Error
}

func (r *UserRepository) FindUserByEmail(email string) (*domainUser.UserModel, error) {
	var pm UserModel
	if err := r.db.Where("email = ?", email).First(&pm).Error; err != nil {
		return nil, err
	}
	d := toDomain(&pm)
	return &d, nil
}

func (r *UserRepository) FindUserByVerificationToken(token string) (*domainUser.UserModel, error) {
	var pm UserModel
	if err := r.db.Where("verification_token = ?", token).First(&pm).Error; err != nil {
		return nil, err
	}
	d := toDomain(&pm)
	return &d, nil
}

func (r *UserRepository) FindUserByPasswordResetToken(token string) (*domainUser.UserModel, error) {
	var pm UserModel
	if err := r.db.Where("password_reset_token = ?", token).First(&pm).Error; err != nil {
		return nil, err
	}
	d := toDomain(&pm)
	return &d, nil
}

func (r *UserRepository) FindByID(id uint) (*domainUser.UserModel, error) {
	var pm UserModel
	if err := r.db.First(&pm, id).Error; err != nil {
		return nil, err
	}
	d := toDomain(&pm)
	return &d, nil
}

func (r *UserRepository) UpdateUser(u *domainUser.UserModel) error {
	pm := toPersistence(u)
	return r.db.Save(&pm).Error
}

func (r *UserRepository) SoftDeleteUnverifiedUsersBefore(cutoff time.Time) error {
	// 未認証(is_verified=false)かつ作成日時(created_at)がcutoff以前のレコードをソフトデリート
	return r.db.
		Where("is_verified = ? AND created_at < ?", false, cutoff).
		Delete(&UserModel{}).
		Error
}

func (r *UserRepository) PermanentlyDeleteUsersBefore(cutoff time.Time) error {
	// 削除日時(deleted_at)がcutoff以前のレコードを完全削除
	return r.db.Unscoped().
		Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoff).
		Delete(&UserModel{}).
		Error
}

// --------------------------------------------------
// toDomain: 永続化モデル -> ドメインモデル
// --------------------------------------------------
func toDomain(pm *UserModel) domainUser.UserModel {
	return domainUser.UserModel{
		ID:                    pm.ID,
		FirstName:             pm.FirstName,
		LastName:              pm.LastName,
		FirstNameKana:         pm.FirstNameKana,
		LastNameKana:          pm.LastNameKana,
		Email:                 pm.Email,
		Password:              pm.Password, // ドメイン側はパスワードをハッシュとして持つ想定
		IsVerified:            pm.IsVerified,
		VerificationToken:     pm.VerificationToken,
		VerificationExpiresAt: pm.VerificationExpiresAt,
		PasswordResetToken:    pm.PasswordResetToken,
		PasswordResetExpires:  pm.PasswordResetExpires,
		SchoolName:            pm.SchoolName,
		Department:            pm.Department,
		Laboratory:            pm.Laboratory,
		GraduationYear:        pm.GraduationYear,
		DesiredJobTypes:       pm.DesiredJobTypes,
		Skills:                pm.Skills,
		SelfIntroduction:      pm.SelfIntroduction,
		ProfileImageURL:       pm.ProfileImageURL,
		CreatedAt:             pm.CreatedAt,
		UpdatedAt:             pm.UpdatedAt,
		DeletedAt:             pm.DeletedAt.Time,
	}
}

// --------------------------------------------------
// toPersistence: ドメインモデル -> 永続化モデル
// --------------------------------------------------
func toPersistence(d *domainUser.UserModel) UserModel {
	return UserModel{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
			DeletedAt: gorm.DeletedAt{Time: d.DeletedAt, Valid: !d.DeletedAt.IsZero()},
		},
		FirstName:             d.FirstName,
		LastName:              d.LastName,
		FirstNameKana:         d.FirstNameKana,
		LastNameKana:          d.LastNameKana,
		Email:                 d.Email,
		Password:              d.Password,
		IsVerified:            d.IsVerified,
		VerificationToken:     d.VerificationToken,
		VerificationExpiresAt: d.VerificationExpiresAt,
		PasswordResetToken:    d.PasswordResetToken,
		PasswordResetExpires:  d.PasswordResetExpires,
		SchoolName:            d.SchoolName,
		Department:            d.Department,
		Laboratory:            d.Laboratory,
		GraduationYear:        d.GraduationYear,
		DesiredJobTypes:       d.DesiredJobTypes,
		Skills:                d.Skills,
		SelfIntroduction:      d.SelfIntroduction,
		ProfileImageURL:       d.ProfileImageURL,
	}
}
