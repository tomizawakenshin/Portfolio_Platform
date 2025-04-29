// backend/domain/user/repository.go
package user

import "time"

// Repository は User エンティティの永続化操作を表すインターフェースです
// サービス層はこのインターフェースだけを依存します
type IUserRepository interface {
	// 新規ユーザーの作成（サインアップ時）
	CreateUser(u *UserModel) error

	// メールアドレスで検索（ログイン時、重複チェック時）
	FindUserByEmail(email string) (*UserModel, error)

	// メール認証用トークンで検索
	FindUserByVerificationToken(token string) (*UserModel, error)

	// パスワードリセット用トークンで検索
	FindUserByPasswordResetToken(token string) (*UserModel, error)

	// 主キーで検索
	FindByID(id uint) (*UserModel, error)

	// プロフィール更新など、変更の保存
	UpdateUser(u *UserModel) error

	// 未認証ユーザーのソフトデリート
	SoftDeleteUnverifiedUsersBefore(cutoff time.Time) error

	// 古い削除済ユーザーの完全削除
	PermanentlyDeleteUsersBefore(cutoff time.Time) error
}
