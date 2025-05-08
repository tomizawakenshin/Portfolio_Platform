// backend/services/auth_service_test.go
package services

import (
	"errors"
	"testing"
	"time"

	domainUser "backend/domain/user"
)

// --- フェイク・リポジトリ ---
type fakeRepo struct {
	// テスト内で記録するフィールド
	createdUser *domainUser.UserModel
	findErr     error
	findUser    *domainUser.UserModel
}

func (f *fakeRepo) FindUserByEmail(email string) (*domainUser.UserModel, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.findUser, nil
}

func (f *fakeRepo) CreateUser(u *domainUser.UserModel) error {
	f.createdUser = u
	return nil
}

// その他のメソッドは今回はダミー実装
func (f *fakeRepo) FindUserByVerificationToken(string) (*domainUser.UserModel, error) {
	return nil, nil
}
func (f *fakeRepo) FindUserByPasswordResetToken(string) (*domainUser.UserModel, error) {
	return nil, nil
}
func (f *fakeRepo) FindByID(uint) (*domainUser.UserModel, error)    { return nil, nil }
func (f *fakeRepo) UpdateUser(*domainUser.UserModel) error          { return nil }
func (f *fakeRepo) SoftDeleteUnverifiedUsersBefore(time.Time) error { return nil }
func (f *fakeRepo) PermanentlyDeleteUsersBefore(time.Time) error    { return nil }

// --- テスト: 新規登録が成功するケース ---
func TestAuthService_SignUp_Success(t *testing.T) {
	// 「未登録」を表すエラー
	repo := &fakeRepo{findErr: errors.New("record not found")}
	svc := NewAuthService(repo)

	token := "verify-token"
	err := svc.SignUp("foo@example.com", "rawpw", token)
	if err != nil {
		t.Fatalf("SignUp failed: %v", err)
	}

	// CreateUser が呼ばれているか
	if repo.createdUser == nil {
		t.Fatal("Expected CreateUser to be called")
	}
	// 引数チェック
	if repo.createdUser.Email != "foo@example.com" {
		t.Errorf("Email mismatch: got %q", repo.createdUser.Email)
	}
	if repo.createdUser.VerificationToken == nil || *repo.createdUser.VerificationToken != token {
		t.Errorf("Token mismatch: got %v", repo.createdUser.VerificationToken)
	}
}

// --- テスト: 既存ユーザーがいるケース ---
func TestAuthService_SignUp_AlreadyExists(t *testing.T) {
	// 既存ユーザーを返す
	existing := &domainUser.UserModel{Email: "foo@example.com"}
	repo := &fakeRepo{findErr: nil, findUser: existing}
	svc := NewAuthService(repo)

	err := svc.SignUp("fo@example.com", "any", "tkn")
	if err == nil || err.Error() != "user already exists" {
		t.Errorf("Expected 'user already exists' error, got %v", err)
	}
}
