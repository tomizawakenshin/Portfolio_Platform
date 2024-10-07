package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUserController interface {
	GetUserInfo(ctx *gin.Context)
}

type UserController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) IUserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) GetUserInfo(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID := user.(*models.User).ID

	// ユーザー情報を取得
	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// ユーザー情報を返す
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
		// 必要に応じて他の情報も追加できます
	})
}
