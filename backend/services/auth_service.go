package services

import (
	"backend/models"
	reposotories "backend/repositories"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	oauth2Google "google.golang.org/api/oauth2/v2"
)

type IAuthService interface {
	SignUp(email string, password string, verificationToken string) error
	Login(email string, password string) (*string, error)
	Logout(ctx *gin.Context) error
	GetUserFromToken(tokenString string) (*models.User, error)
	VerifyUser(token string) (*models.User, error)
	CreateToken(userId uint, email string) (*string, error)
	SoftDeleteUnverifiedUsers() error
	PermanentlyDeleteUsers() error
	FindOrCreateUserByGoogle(userinfo *oauth2Google.Userinfo) (*models.User, error)
}

type AuthService struct {
	repository reposotories.IAuthRepository
}

func NewAuthService(repository reposotories.IAuthRepository) IAuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) SignUp(email string, password string, verificationToken string) error {
	// ユーザーが既に存在するか確認
	_, err := s.repository.FindUser(email)
	if err == nil {
		// ユーザーが存在する場合はエラーを返す
		return errors.New("user already exists")
	}

	// 新しいユーザーを作成
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:             email,
		Password:          string(hashedPassword),
		IsVerified:        false,
		VerificationToken: verificationToken,
	}

	return s.repository.CreateUser(user)
}

func (s *AuthService) VerifyUser(token string) (*models.User, error) {
	user, err := s.repository.FindUserByToken(token)
	if err != nil {
		return nil, err
	}

	// isVerifiedをtrueに更新
	user.IsVerified = true
	user.VerificationToken = "" // トークンはクリアする
	if err := s.repository.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email string, password string) (*string, error) {
	foundUser, err := s.repository.FindUser(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	jwtToken, err := s.CreateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}

func (s *AuthService) CreateToken(userId uint, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
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

		user, err = s.repository.FindUser(claims["email"].(string))
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
	user, err := s.repository.FindUser(userinfo.Email)
	if err != nil {
		if err.Error() != "user not found" {
			return nil, err
		}
		// ユーザーが存在しない場合、新規作成
		user = &models.User{
			Email:      userinfo.Email,
			FirstName:  "",
			LastName:   "",
			IsVerified: true,
			// パスワードは空またはランダムな値を設定
			Password: "",
		}
		if err := s.repository.CreateUser(*user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
