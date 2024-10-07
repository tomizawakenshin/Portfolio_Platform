package main

import (
	"backend/controllers"
	"backend/infra"
	"backend/middlewares"
	reposotories "backend/repositories"
	"backend/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {

	authRepository := reposotories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	emailService := services.NewEmailService()
	authController := controllers.NewAuthController(authService, emailService)

	userRepository := reposotories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
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

	//user認証関連のエンドポイント
	authRouter := r.Group("/auth")
	authRouter.POST("/signup", authController.SignUp)
	authRouter.POST("/login", authController.Login)
	authRouter.GET("/verify", authController.VerifyAccount)

	//user情報関連のエンドポイント
	userRouterWithAuth := r.Group("/user", middlewares.AuthMiddleware(authService))
	userRouterWithAuth.GET("/GetInfo", userController.GetUserInfo)

	return r

}

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	r := setupRouter(db)
	r.Run("localhost:8080") // 0.0.0.0:8080 でサーバーを立てます。
}
