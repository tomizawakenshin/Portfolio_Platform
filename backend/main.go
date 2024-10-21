package main

import (
	"backend/controllers"
	"backend/infra"
	"backend/middlewares"
	reposotories "backend/repositories"
	"backend/services"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB, authService services.IAuthService) *gin.Engine {

	authRepository := reposotories.NewAuthRepository(db)
	// authService := services.NewAuthService(authRepository)
	emailService := services.NewEmailService()
	authController := controllers.NewAuthController(authService, emailService)

	userRepository := reposotories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, authRepository)
	userController := controllers.NewUserController(userService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.vercel.app", "http://localhost:8025/#"}, // フロントエンドのドメインを許可
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                                              // 許可するHTTPメソッド
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},                                              // 許可するリクエストヘッダー
		ExposeHeaders:    []string{"Content-Length"},                                                                       // クライアントに公開するレスポンスヘッダー
		AllowCredentials: true,                                                                                             // 認証情報（クッキーなど）の送信を許可
		MaxAge:           48 * time.Hour,                                                                                   // プリフライトリクエストのキャッシュ時間
	}))

	//user認証のエンドポイント
	authRouter := r.Group("/auth")
	authRouter.POST("/signup", authController.SignUp)
	authRouter.POST("/login", authController.Login)
	authRouter.GET("/verify", authController.VerifyAccount)

	// Google OAuth 2.0認証のエンドポイント
	authRouter.GET("/google/login", authController.GoogleLogin)
	authRouter.GET("/google/callback", authController.GoogleCallback)

	//ログアウトのエンドポイント
	authRouterWithAuth := r.Group("/auth", middlewares.AuthMiddleware(authService))
	authRouterWithAuth.POST("/logout", authController.Logout)

	//Cookieの存在の確認用のエンドポイント
	authRouter.GET("/check", authController.CheckAuth)

	//パスワードリセットのエンドポイント
	authRouter.POST("/RequestPasswordReset", authController.RequestPasswordReset)
	authRouter.POST("/CheckResetToken", authController.CheckResetToken)
	authRouter.POST("/ResetPassword", authController.ResetPassword)

	//user情報関連のエンドポイント
	userRouterWithAuth := r.Group("/user", middlewares.AuthMiddleware(authService))
	userRouterWithAuth.GET("/GetInfo", userController.GetUserInfo)
	userRouterWithAuth.PUT("/UpdateMinimumUserInfo", userController.UpdateMinimumUserInfo)

	return r

}

func startSoftDeleteJob(authService services.IAuthService) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			err := authService.SoftDeleteUnverifiedUsers()
			if err != nil {
				log.Printf("Error soft-deleting unverified users: %v", err)
			} else {
				log.Println("Soft-delete job executed successfully")
			}
		}
	}()
}

func startPermanentDeletionJob(authService services.IAuthService) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			err := authService.PermanentlyDeleteUsers()
			if err != nil {
				log.Printf("Error permanently deleting users: %v", err)
			} else {
				log.Println("Permanent deletion job executed successfully")
			}
		}
	}()
}

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	authRepository := reposotories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)

	// クリーンアップジョブの開始
	startSoftDeleteJob(authService)
	startPermanentDeletionJob(authService)

	r := setupRouter(db, authService)
	r.Run("localhost:8080") // 0.0.0.0:8080 でサーバーを立てます。
}
