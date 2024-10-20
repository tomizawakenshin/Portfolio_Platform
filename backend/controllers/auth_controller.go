package controllers

import (
	"backend/dto"
	"backend/services"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"
	"gorm.io/gorm"
)

type IAuthController interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	VerifyAccount(ctx *gin.Context)
	GoogleLogin(ctx *gin.Context)
	GoogleCallback(ctx *gin.Context)
	CheckAuth(ctx *gin.Context)
	RequestPasswordReset(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	CheckResetToken(ctx *gin.Context)
}

type AuthController struct {
	services          services.IAuthService
	emailService      services.IEmailService
	googleOauthConfig *oauth2.Config
}

func NewAuthController(service services.IAuthService, emailService services.IEmailService) IAuthController {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"openid", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return &AuthController{
		services:          service,
		emailService:      emailService,
		googleOauthConfig: googleOauthConfig,
	}
}

func (c *AuthController) SignUp(ctx *gin.Context) {
	var input dto.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 本登録用のトークンを生成
	verificationToken := generateVerificationToken()

	err := c.services.SignUp(input.Email, input.Password, verificationToken)
	if err != nil {
		if err.Error() == "user already exists" {
			// ユーザーが既に存在する場合はログイン処理に移行
			jwtToken, tokenExpiry, err := c.services.Login(input.Email, input.Password, false)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			// 【変更点①】既存ユーザーでログインした場合、ステータスコードを http.StatusOK（200）に設定
			ctx.SetCookie("jwt-token", *jwtToken, int(tokenExpiry.Seconds()), "/", "localhost", false, true)
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Logged in successfully",
				"token":   jwtToken})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Temporary User"})
		return
	}

	// メール送信
	if err := c.emailService.SendRegistrationEmail(input.Email, verificationToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	// 【変更点②】新規ユーザー作成時、ステータスコードを http.StatusCreated（201）に設定
	ctx.JSON(http.StatusCreated, gin.H{"message": "仮登録が完了しました。メールを確認してください。"})
}

func generateVerificationToken() string {
	// 16バイトのランダムなトークンを生成
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	// トークンを16進数にエンコード
	return hex.EncodeToString(token)
}

func (c *AuthController) VerifyAccount(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	user, err := c.services.VerifyUser(token)
	if err != nil {
		switch err {
		case services.ErrUserAlreadyVerified:
			ctx.Redirect(http.StatusFound, "http://localhost:3000/auth")
		case services.ErrVerificationTokenExpired:
			ctx.Redirect(http.StatusFound, "http://localhost:3000/auth")
		case services.ErrUserNotFound:
			ctx.Redirect(http.StatusFound, "http://localhost:3000/auth")
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := c.emailService.SendWelcomeEmail(user.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	jwtToken, tokenExpiry, err := c.services.CreateToken(user.ID, user.Email, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT token"})
		return
	}

	ctx.SetCookie("jwt-token", *jwtToken, int(tokenExpiry.Seconds()), "/", "localhost", false, true)

	ctx.Redirect(http.StatusFound, "http://localhost:3000/home")
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, tokenExpiry, err := c.services.Login(input.Email, input.Password, input.RememberMe)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("jwt-token", *jwtToken, int(tokenExpiry.Seconds()), "/", "localhost", false, true)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	err := c.services.Logout(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// クッキーを削除
	ctx.SetCookie("jwt-token", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "ログアウトしました。"})
}

func (c *AuthController) GoogleLogin(ctx *gin.Context) {
	rememberMe := ctx.Query("rememberMe")
	state := "state-token"

	// rememberMeフラグをstateに含める
	if rememberMe == "true" {
		state += "|rememberMe"
	}

	url := c.googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *AuthController) GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	rememberMe := false

	if strings.Contains(state, "|rememberMe") {
		rememberMe = true
		state = strings.Replace(state, "|rememberMe", "", 1)
	}

	if state != "state-token" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "State token does not match"})
		return
	}

	code := ctx.Query("code")
	token, err := c.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// ユーザー情報を取得
	client := c.googleOauthConfig.Client(context.Background(), token)
	oauth2Service, err := oauth2v2.New(client)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create oauth2 service"})
		return
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// ユーザーをデータベースに作成または取得
	user, err := c.services.FindOrCreateUserByGoogle(userinfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create user"})
		return
	}

	if err := c.emailService.SendWelcomeEmail(user.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	// JWTトークンを作成
	jwtToken, tokenExpiry, err := c.services.CreateToken(user.ID, user.Email, rememberMe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT token"})
		return
	}

	// クッキーにトークンを設定
	ctx.SetCookie("jwt-token", *jwtToken, int(tokenExpiry.Seconds()), "/", "localhost", false, true)

	// フロントエンドにリダイレクト
	ctx.Redirect(http.StatusFound, "http://localhost:3000/home")
}

func (c *AuthController) CheckAuth(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("jwt-token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// トークンを検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// トークンの署名方法を検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Authorized"})
}

func (c *AuthController) RequestPasswordReset(ctx *gin.Context) {
	var input dto.PasswordResetRequestInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resetToken, err := c.services.GeneratePasswordResetToken(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ユーザーが存在しない場合
			ctx.JSON(http.StatusNotFound, gin.H{"error": "そのアカウントは無効です。"})
		} else {
			// その他のエラー
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました。"})
		}
		return
	}

	// パスワードリセットメールを送信
	err = c.emailService.SendPasswordResetEmail(input.Email, resetToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードリセットメールの送信に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "メールアドレスにパスワードリセットのリンクを送信しました。"})
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var input dto.PasswordResetInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです。"})
		return
	}

	user, err := c.services.ValidatePasswordResetToken(input.Token)
	if err != nil {
		if err.Error() == "reset token has expired" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "トークンの有効期限が切れています。"})
		} else if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "無効なトークンです。"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました。"})
		}
		return
	}

	err = c.services.UpdatePassword(user, input.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの更新に失敗しました"})
		return
	}

	// パスワードリセット完了メールを送信
	if err := c.emailService.SendPasswordResetConfirmationEmail(user.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードリセット完了メールの送信に失敗しました"})
		return
	}

	// オプション：パスワードリセット後にユーザーをログインさせる
	jwtToken, tokenExpiry, err := c.services.CreateToken(user.ID, user.Email, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWTトークンの作成に失敗しました"})
		return
	}

	ctx.SetCookie("jwt-token", *jwtToken, int(tokenExpiry.Seconds()), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "パスワードがリセットされ、ログインしました。"})
}

func (c *AuthController) CheckResetToken(ctx *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "トークンが必要です。"})
		return
	}

	_, err := c.services.ValidatePasswordResetToken(input.Token)
	if err != nil {
		// エラーの種類に応じて適切なステータスコードを返す
		if err.Error() == "reset token has expired" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "トークンの有効期限が切れています。"})
		} else if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "無効なトークンです。"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました。"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "トークンは有効です。"})
}
