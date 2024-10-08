package controllers

import (
	"backend/dto"
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUserController interface {
	GetUserInfo(ctx *gin.Context)
	UpdateMinimumUserInfo(ctx *gin.Context)
}

type UserController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) IUserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	// コンテキストからユーザーIDを取得
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID := user.(*models.User).ID

	// ユーザー情報を取得
	user, err := c.userService.GetUserByID(userID)
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

func (c *UserController) UpdateMinimumUserInfo(ctx *gin.Context) {
	// コンテキストからユーザーを取得
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := user.(*models.User)

	// リクエストボディをバインド
	var input dto.MinimumUserInfoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// サービスを呼び出してユーザー情報を更新
	if err := c.userService.UpdateMinimumUserInfo(currentUser.ID, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		return
	}

	// 成功レスポンスを返す
	ctx.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
}
