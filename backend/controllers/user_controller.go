// controllers/user_controller.go

package controllers

import (
	domainUser "backend/domain/user"
	"backend/dto"
	"backend/services"
	"mime/multipart"
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

	userID := user.(*domainUser.UserModel).ID

	// ユーザー情報を取得
	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// ユーザー情報を返す
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (c *UserController) UpdateMinimumUserInfo(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := user.(*domainUser.UserModel)

	var input dto.MinimumUserInfoInput
	var fileHeaders []*multipart.FileHeader

	// リクエストのContent-Typeを取得
	contentType := ctx.Request.Header.Get("Content-Type")

	// JSONリクエストの場合の処理
	if contentType == "application/json" || contentType == "application/json; charset=utf-8" {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "details": err.Error()})
			return
		}
	} else {
		// multipart/form-dataの場合の処理
		if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form", "details": err.Error()})
			return
		}
		if err := ctx.ShouldBind(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
			return
		}
		form, _ := ctx.MultipartForm()
		fileHeaders = form.File["profileImage"]
	}

	// サービス層へ
	updatedUser, err := c.userService.UpdateMinimumUserInfo(currentUser.ID, input, fileHeaders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User info updated",
		"user":    updatedUser,
	})
}
