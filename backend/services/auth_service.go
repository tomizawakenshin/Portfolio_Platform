package services

import (
	"backend/models"
	reposotories "backend/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignUp(email string, password string, verificationToken string) error
	Login(email string, password string) (*string, error)
	// GetUserFromToken(tokenString string) (*models.User, error)
	VerifyUser(token string) error
}

type AuthService struct {
	repository reposotories.IAuthRepository
}

func NewAuthService(repository reposotories.IAuthRepository) IAuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) SignUp(email string, password string, verificationToken string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:             email,
		Password:          string(hashedPassword),
		IsVerified:        false,             // 仮登録として作成
		VerificationToken: verificationToken, // 本登録用トークンを保存
	}

	return s.repository.CreateUser(user)
}

func (s *AuthService) VerifyUser(token string) error {
	user, err := s.repository.FindUserByToken(token)
	if err != nil {
		return err
	}

	// isVerifiedをtrueに更新
	user.IsVerified = true
	user.VerificationToken = "" // トークンはクリアする
	return s.repository.UpdateUser(user)
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

	token, err := CreateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func CreateToken(userId uint, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 5).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
