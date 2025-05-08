// backend/domain/user/entity_test.go
package user

import (
	"strings"
	"testing"
	"time"
)

func TestNewUser_Validation(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		wantErrPart string // "" のときはエラーなし
	}{
		{
			name:        "empty email",
			email:       "",
			password:    "hashedpw",
			wantErrPart: "メールアドレスは必須です",
		},
		{
			name:        "invalid email format",
			email:       "not-an-email",
			password:    "hashedpw",
			wantErrPart: "無効なメールアドレス形式です",
		},
		{
			name:        "empty password",
			email:       "foo@example.com",
			password:    "",
			wantErrPart: "パスワードは必須です",
		},
		{
			name:        "valid inputs",
			email:       "foo@example.com",
			password:    "hashedpw",
			wantErrPart: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser(tt.email, tt.password)

			if tt.wantErrPart != "" {
				// エラーが返ってくるパターン
				if err == nil {
					t.Fatalf("expected error containing %q, but got nil", tt.wantErrPart)
				}
				if !strings.Contains(err.Error(), tt.wantErrPart) {
					t.Errorf("expected error containing %q, but got %q", tt.wantErrPart, err.Error())
				}
				if u != nil {
					t.Errorf("expected user to be nil on error, but got %+v", u)
				}
				return
			}

			// 正常系: err は nil、返ってきた UserModel の検証
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if u.Email != tt.email {
				t.Errorf("Email = %q; want %q", u.Email, tt.email)
			}
			if u.Password == nil || *u.Password != tt.password {
				t.Errorf("Password = %v; want pointer to %q", u.Password, tt.password)
			}
			// CreatedAt/UpdatedAt がセットされているか
			if u.CreatedAt.IsZero() || u.UpdatedAt.IsZero() {
				t.Error("CreatedAt or UpdatedAt is zero")
			}
			// DeletedAt は初期値 (IsZero が true)
			if !u.DeletedAt.IsZero() {
				t.Errorf("DeletedAt should be zero on NewUser, got %v", u.DeletedAt)
			}
			// VerificationToken / VerificationExpiresAt は空
			if u.VerificationToken != nil {
				t.Errorf("VerificationToken should be nil on NewUser, got %v", u.VerificationToken)
			}
			if !u.VerificationExpiresAt.IsZero() {
				t.Errorf("VerificationExpiresAt should be zero on NewUser, got %v", u.VerificationExpiresAt)
			}
			// PasswordResetToken / PasswordResetExpires は空
			if u.PasswordResetToken != "" {
				t.Errorf("PasswordResetToken should be empty, got %q", u.PasswordResetToken)
			}
			if !u.PasswordResetExpires.IsZero() {
				t.Errorf("PasswordResetExpires should be zero, got %v", u.PasswordResetExpires)
			}
			// そのほかのフィールド（プロフィール等）は空文字 or 空スライス
			if u.FirstName != "" || u.LastName != "" {
				t.Errorf("expected empty names, got %q %q", u.FirstName, u.LastName)
			}
			if len(u.DesiredJobTypes) != 0 || len(u.Skills) != 0 {
				t.Errorf("expected empty slices, got DesiredJobTypes=%v, Skills=%v", u.DesiredJobTypes, u.Skills)
			}
			// 最終的に CreatedAt と UpdatedAt の差がごく短いことも確認できます
			if u.UpdatedAt.Sub(u.CreatedAt) > time.Second {
				t.Errorf("UpdatedAt should be almost equal CreatedAt, diff = %v", u.UpdatedAt.Sub(u.CreatedAt))
			}
		})
	}
}
