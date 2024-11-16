package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IPortfolioController interface {
	CreatePost(ctx *gin.Context)
	GetPostsByUserID(ctx *gin.Context)
	GetAllPosts(ctx *gin.Context)
}

type PortfolioController struct {
	portfolioService services.IPortfolioService
}

func NewPortfolioController(portfolioService services.IPortfolioService) IPortfolioController {
	return &PortfolioController{portfolioService: portfolioService}
}

func (c *PortfolioController) CreatePost(ctx *gin.Context) {
	// ユーザー認証情報を取得
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := user.(*models.User)

	// フォームの最大サイズを設定（例: 32MB）
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form", "details": err.Error()})
		return
	}

	// フォームデータをサービス層に渡す
	if err := c.portfolioService.CreatePost(ctx, currentUser.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post", "details": err.Error()})
		return
	}

	// 成功レスポンスを返す
	ctx.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

func (c *PortfolioController) GetPostsByUserID(ctx *gin.Context) {
	// ユーザー認証情報を取得
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := user.(*models.User)

	// サービスを呼び出して投稿一覧を取得
	posts, err := c.portfolioService.GetPostsByUserID(currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	// 投稿一覧を返す
	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (c *PortfolioController) GetAllPosts(ctx *gin.Context) {
	portfolio, err := c.portfolioService.GetAllPosts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}
