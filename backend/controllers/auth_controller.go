package controllers

import (
	"backend/dto"
	"backend/services"
	"crypto/rand"
	"encoding/hex"

	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	VerifyAccount(ctx *gin.Context)
}

type AuthController struct {
	services     services.IAuthService
	emailService services.IEmailService
}

func NewAuthController(service services.IAuthService, emailService services.IEmailService) IAuthController {
	return &AuthController{
		services:     service,
		emailService: emailService,
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Temporary User"})
		return
	}

	// 本登録リンクを生成（トークン付き）
	// verificationLink := fmt.Sprintf("https://your-frontend.vercel.app/verify?token=%s", verificationToken)

	// メール送信
	if err := c.emailService.SendRegistrationEmail(input.Email, verificationToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

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

	// ctx.JSON(http.StatusOK, gin.H{"message": "アカウントが有効化されました。"})
	ctx.Redirect(http.StatusFound, "http://localhost:3000/home")
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.services.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
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
