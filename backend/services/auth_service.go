package services

import (
	"backend/models"
	"backend/repositories"
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
	GetUserFromToken(tokenString string) (*models.User, error)
	VerifyUser(token string) (*models.User, error)
	CreateToken(userId uint, email string, rememberMe bool) (*string, time.Duration, error)
	SoftDeleteUnverifiedUsers() error
	PermanentlyDeleteUsers() error
	FindOrCreateUserByGoogle(userinfo *oauth2Google.Userinfo) (*models.User, error)
	GeneratePasswordResetToken(email string) (string, error)
	ValidatePasswordResetToken(token string) (*models.User, error)
	UpdatePassword(user *models.User, newPassword string) error
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
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
	user := models.User{
		Email:                 email,
		Password:              &hashedPasswordStr,
		IsVerified:            false,
		VerificationToken:     &verificationToken,
		VerificationExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	return s.repository.CreateUser(user)
}

func (s *AuthService) VerifyUser(token string) (*models.User, error) {
	user, err := s.repository.FindUserByVerificationToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// トークンの有効期限を確認
	if time.Now().After(user.VerificationExpiresAt) {
		return nil, ErrVerificationTokenExpired
	}

	// 既に本登録済みの場合
	if user.IsVerified {
		return nil, ErrUserAlreadyVerified
	}

	// isVerifiedをtrueに更新
	user.IsVerified = true
	user.VerificationToken = nil             // トークンはクリアする
	user.VerificationExpiresAt = time.Time{} // 有効期限をクリア

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

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	var user *models.User
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

func (s *AuthService) FindOrCreateUserByGoogle(userinfo *oauth2Google.Userinfo) (*models.User, error) {
	// メールアドレスでユーザーを検索
	user, err := s.repository.FindUserByEmail(userinfo.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ユーザーが存在しない場合、新規作成
			user = &models.User{
				Email:      userinfo.Email,
				FirstName:  "",
				LastName:   "",
				IsVerified: true,
				Password:   nil,
			}
			if err := s.repository.CreateUser(*user); err != nil {
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

	// トークンと有効期限を設定（例：1時間後に有効期限切れ）
	user.PasswordResetToken = resetToken
	user.PasswordResetExpires = time.Now().Add(time.Hour)

	// 新しいトークンでユーザーを保存
	err = s.repository.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *AuthService) ValidatePasswordResetToken(token string) (*models.User, error) {
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

func (s *AuthService) UpdatePassword(user *models.User, newPassword string) error {
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
