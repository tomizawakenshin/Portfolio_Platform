package controllers

import (
	"backend/dto"
	"backend/services"
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"

	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"
)

type IAuthController interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	VerifyAccount(ctx *gin.Context)
	GoogleLogin(ctx *gin.Context)
	GoogleCallback(ctx *gin.Context)
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
			jwtToken, err := c.services.Login(input.Email, input.Password)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			// 【変更点①】既存ユーザーでログインした場合、ステータスコードを http.StatusOK（200）に設定
			ctx.SetCookie("jwt-token", *jwtToken, 3600*24, "/", "localhost", false, true)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := c.services.CreateToken(user.ID, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT token"})
		return
	}

	ctx.SetCookie("jwt-token", *jwtToken, 3600*24, "/", "localhost", false, true)

	ctx.Redirect(http.StatusFound, "http://localhost:3000/home")
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := c.services.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("jwt-token", *jwtToken, 3600*24, "/", "localhost", false, true)
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
	url := c.googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *AuthController) GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
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

	// JWTトークンを作成
	jwtToken, err := c.services.CreateToken(user.ID, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT token"})
		return
	}

	// クッキーにトークンを設定
	ctx.SetCookie("jwt-token", *jwtToken, 3600*24, "/", "localhost", false, true)

	// フロントエンドにリダイレクト
	ctx.Redirect(http.StatusFound, "http://localhost:3000/home")
}
