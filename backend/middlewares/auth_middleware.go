package middlewares

import (
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("jwt-token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// トークンの検証とユーザーの取得
		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
