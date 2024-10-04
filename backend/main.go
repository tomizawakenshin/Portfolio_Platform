package main

import (
	"backend/infra"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong!",
		})
	})
	r.Run("localhost:8080") // 0.0.0.0:8080 でサーバーを立てます。
}
