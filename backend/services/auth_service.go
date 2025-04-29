package services

import (
	domainUser "backend/domain/user"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	oauth2Google "google.golang.org/api/oauth2/v2"
)

type IAuthService interface {
	SignUp(email string, password string, verificationToken string) error
	Login(email string, password string, rememberMe bool) (*string, time.Duration, error)
	Logout(ctx *gin.Context) error
	GetUserFromToken(tokenString string) (*domainUser.UserModel, error)
	VerifyUser(token string) (*domainUser.UserModel, error)
	CreateToken(userId uint, email string, rememberMe bool) (*string, time.Duration, error)
	SoftDeleteUnverifiedUsers() error
	PermanentlyDeleteUsers() error
	FindOrCreateUserByGoogle(userinfo *oauth2Google.Userinfo) (*domainUser.UserModel, error)
	GeneratePasswordResetToken(email string) (string, error)
	ValidatePasswordResetToken(token string) (*domainUser.UserModel, error)
	UpdatePassword(user *domainUser.UserModel, newPassword string) error
}

type AuthService struct {
	// repository repositories.IAuthRepository
	repository domainUser.IUserRepository
}

func NewAuthService(repository domainUser.IUserRepository) IAuthService {
	return &AuthService{repository: repository}
}

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyVerified      = errors.New("user already verified")
	ErrVerificationTokenExpired = errors.New("verification token has expired")
)

func (s *AuthService) SignUp(email string, password string, verificationToken string) error {
	// ユーザーが既に存在するか確認
	_, err := s.repository.FindUserByEmail(email)
	if err == nil {
		// ユーザーが存在する場合はエラーを返す
		return errors.New("user already exists")
	}

	// 新しいユーザーを作成
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedPasswordStr := string(hashedPassword)
	user, err := domainUser.NewUser(email, hashedPasswordStr)
	if err != nil {
		return err
	}
	user.VerificationToken = &verificationToken
	user.VerificationExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	return s.repository.CreateUser(user)
}

func (s *AuthService) VerifyUser(token string) (*domainUser.UserModel, error) {
	user, err := s.repository.FindUserByVerificationToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// ドメイン層の検証ロジックを使う
	if err := user.VerifyEmail(token); err != nil {
		// 既に認証済み or トークン不正 などはここでエラーとして返してくれる
		return nil, err
	}

	if err := s.repository.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email string, password string, rememberMe bool) (*string, time.Duration, error) {
	foundUser, err := s.repository.FindUserByEmail(email)
	if err != nil {
		return nil, 0, err
	}

	if foundUser.Password == nil {
		return nil, 0, errors.New("パスワードが設定されていません")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*foundUser.Password), []byte(password))
	if err != nil {
		return nil, 0, err
	}

	jwtToken, tokenExpiry, err := s.CreateToken(foundUser.ID, foundUser.Email, rememberMe)
	if err != nil {
		return nil, 0, err
	}

	return jwtToken, tokenExpiry, nil
}

func (s *AuthService) CreateToken(userId uint, email string, rememberMe bool) (*string, time.Duration, error) {
	var tokenExpiry time.Duration
	if rememberMe {
		tokenExpiry = time.Hour * 24 * 14 // 14日間
	} else {
		tokenExpiry = time.Hour * 1 // 1時間
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(tokenExpiry).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, 0, err
	}

	return &tokenString, tokenExpiry, nil
}

func (s *AuthService) GetUserFromToken(tokenString string) (*domainUser.UserModel, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	var user *domainUser.UserModel
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, jwt.ErrTokenExpired
		}

		user, err = s.repository.FindUserByEmail(claims["email"].(string))
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (s *AuthService) Logout(ctx *gin.Context) error {
	// 必要に応じてトークンのブラックリスト化などを実装
	return nil
}

// ソフトデリートを行うメソッドを追加
func (s *AuthService) SoftDeleteUnverifiedUsers() error {
	cutoffTime := time.Now().UTC().Add(-168 * time.Hour) // 本登録の有効期間は7日間
	return s.repository.SoftDeleteUnverifiedUsersBefore(cutoffTime)
}

// ハードデリートを行うメソッドを追加
func (s *AuthService) PermanentlyDeleteUsers() error {
	cutoffTime := time.Now().UTC().Add(-552 * time.Hour) // ソフトデリートされてからハードデリートされるまでは3週間
	return s.repository.PermanentlyDeleteUsersBefore(cutoffTime)
}

func (s *AuthService) FindOrCreateUserByGoogle(userinfo *oauth2Google.Userinfo) (*domainUser.UserModel, error) {
	// メールアドレスでユーザーを検索
	user, err := s.repository.FindUserByEmail(userinfo.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ユーザーが存在しない場合、新規作成
			user = &domainUser.UserModel{
				Email:      userinfo.Email,
				FirstName:  "",
				LastName:   "",
				IsVerified: true,
				Password:   nil,
			}
			if err := s.repository.CreateUser(user); err != nil {
				return nil, err
			}
		} else {
			fmt.Printf("Error finding user by email: %v", err)
			return nil, err
		}
	}

	return user, nil
}

func (s *AuthService) GeneratePasswordResetToken(email string) (string, error) {
	user, err := s.repository.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	// ランダムなトークンを生成
	tokenBytes := make([]byte, 16)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	resetToken := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(time.Hour)

	user.RequestPasswordReset(resetToken, expiresAt)

	// トークンと有効期限を設定（例：1時間後に有効期限切れ）
	// user.PasswordResetToken = resetToken
	// user.PasswordResetExpires = time.Now().Add(time.Hour)

	// 新しいトークンでユーザーを保存
	err = s.repository.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *AuthService) ValidatePasswordResetToken(token string) (*domainUser.UserModel, error) {
	user, err := s.repository.FindUserByPasswordResetToken(token)
	if err != nil {
		return nil, err
	}

	// トークンが有効期限切れか確認
	if time.Now().After(user.PasswordResetExpires) {
		return nil, errors.New("reset token has expired")
	}

	return user, nil
}

func (s *AuthService) UpdatePassword(user *domainUser.UserModel, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedPasswordStr := string(hashedPassword)
	user.Password = &hashedPasswordStr
	// リセットトークンと有効期限をクリア
	user.PasswordResetToken = ""
	user.PasswordResetExpires = time.Time{}

	return s.repository.UpdateUser(user)
}
