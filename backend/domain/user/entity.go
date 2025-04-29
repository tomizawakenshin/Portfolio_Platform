// backend/domain/user/entity.go
package user

import (
	"fmt"
	"net/mail"
	"time"
)

// User はユーザーに関するドメインエンティティです。
// フィールドは DB のスキーマに依存せず、ビジネスロジックに沿った形で保持します。
type UserModel struct {
	ID                    uint      // 永続化後に DB から設定される
	Email                 string    // ログインに使うメールアドレス
	Password              *string   // ハッシュ化したパスワード
	IsVerified            bool      // メール認証済みフラグ
	VerificationToken     *string   // メール認証用トークン
	VerificationExpiresAt time.Time // トークンの有効期限
	PasswordResetToken    string    // パスワードリセット用トークン
	PasswordResetExpires  time.Time // リセットトークンの有効期限

	// プロフィール情報（最低限必要な情報）
	FirstName        string
	LastName         string
	FirstNameKana    string
	LastNameKana     string
	ProfileImageURL  string
	SelfIntroduction string

	SchoolName      string
	Department      string
	Laboratory      string
	GraduationYear  string // GORMモデルが文字列なので、ドメインも string に合わせる
	DesiredJobTypes []string
	Skills          []string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// NewUser は新規登録時に呼ぶファクトリメソッドです。
// ・email は必須かつ正しい形式であること
// ・password はすでにハッシュ化済みの文字列として受け取り、非空であること
func NewUser(email, password string) (*UserModel, error) {
	// email 空チェック & フォーマットチェック
	if email == "" {
		return nil, fmt.Errorf("メールアドレスは必須です")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("無効なメールアドレス形式です")
	}

	// password はすでにハッシュ化された文字列として受け取り、空文字は許さない
	if password == "" {
		return nil, fmt.Errorf("パスワードは必須です")
	}

	now := time.Now()
	pw := password
	return &UserModel{
		Email:                 email,
		Password:              &pw,
		IsVerified:            false,
		VerificationToken:     nil,
		VerificationExpiresAt: time.Time{},
		PasswordResetToken:    "",
		PasswordResetExpires:  time.Time{},

		FirstName:       "",
		LastName:        "",
		FirstNameKana:   "",
		LastNameKana:    "",
		ProfileImageURL: "",

		SchoolName:      "",
		Department:      "",
		Laboratory:      "",
		GraduationYear:  "",
		DesiredJobTypes: []string{},
		Skills:          []string{},

		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// VerifyEmail はトークンを受け取ってメール認証を完了させる振る舞い
func (u *UserModel) VerifyEmail(token string) error {
	if time.Now().After(u.VerificationExpiresAt) {
		return fmt.Errorf("認証トークンの有効期限が切れています")
	}

	if u.IsVerified {
		return fmt.Errorf("すでに認証済みです")
	}
	if u.VerificationToken == nil || token != *u.VerificationToken {
		return fmt.Errorf("認証トークンが不正です")
	}
	u.IsVerified = true
	u.VerificationToken = nil
	u.UpdatedAt = time.Now()
	return nil
}

// RequestPasswordReset はパスワードリセットトークンと有効期限をセットする振る舞い
func (u *UserModel) RequestPasswordReset(token string, expiresAt time.Time) {
	t := token
	u.PasswordResetToken = t
	u.PasswordResetExpires = expiresAt
	u.UpdatedAt = time.Now()
}

// ResetPassword はトークン検証後に新しいハッシュパスワードをセットする振る舞い
func (u *UserModel) ResetPassword(token, newHash string) error {
	if u.PasswordResetToken == "" || token != u.PasswordResetToken {
		return fmt.Errorf("パスワードリセットトークンが不正です")
	}
	if newHash == "" {
		return fmt.Errorf("新しいパスワードは必須です")
	}
	ph := newHash
	u.Password = &ph
	u.PasswordResetToken = ""
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile はプロフィール情報を一括で更新する振る舞い
func (u *UserModel) UpdateProfile(
	firstName, lastName, firstNameKana, lastNameKana,
	schoolName, department, laboratory string,
	graduationYear string,
	desiredJobTypes, skills []string,
) error {
	if firstName == "" || lastName == "" {
		return fmt.Errorf("氏名（漢字）は必須です")
	}
	if firstNameKana == "" || lastNameKana == "" {
		return fmt.Errorf("氏名（カナ）は必須です")
	}
	if graduationYear == "" {
		return fmt.Errorf("卒業年が不正です")
	}
	u.FirstName = firstName
	u.LastName = lastName
	u.FirstNameKana = firstNameKana
	u.LastNameKana = lastNameKana
	u.SchoolName = schoolName
	u.Department = department
	u.Laboratory = laboratory
	u.GraduationYear = graduationYear
	u.DesiredJobTypes = desiredJobTypes
	u.Skills = skills
	u.UpdatedAt = time.Now()
	return nil
}
