package main

import (
	"backend/controllers"
	"backend/infra"
	reposotories "backend/repositories"
	"backend/services"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {

	authRepository := reposotories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.vercel.app"}, // フロントエンドのドメインを許可
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                   // 許可するHTTPメソッド
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},                   // 許可するリクエストヘッダー
		ExposeHeaders:    []string{"Content-Length"},                                            // クライアントに公開するレスポンスヘッダー
		AllowCredentials: true,                                                                  // 認証情報（クッキーなど）の送信を許可
		MaxAge:           48 * time.Hour,                                                        // プリフライトリクエストのキャッシュ時間
	}))

	r.OPTIONS("/hanabi/getAll", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(http.StatusNoContent) // 204 No Contentを返す
	})

	//user認証関連のエンドポイント
	authRouter := r.Group("/auth")
	authRouter.POST("/signup", authController.SignUp)
	authRouter.POST("/login", authController.Login)

	return r

}

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	r := setupRouter(db)
	r.Run("localhost:8080") // 0.0.0.0:8080 でサーバーを立てます。
}
