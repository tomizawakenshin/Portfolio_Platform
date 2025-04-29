package controllers

import (
	domainUser "backend/domain/user"
	"backend/dto"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IPortfolioController interface {
	CreatePost(ctx *gin.Context)
	GetPostsByUserID(ctx *gin.Context)
	GetAllPosts(ctx *gin.Context)
	GetPostByID(ctx *gin.Context)
}

type PortfolioController struct {
	portfolioService services.IPortfolioService
}

func NewPortfolioController(portfolioService services.IPortfolioService) IPortfolioController {
	return &PortfolioController{portfolioService: portfolioService}
}

func (c *PortfolioController) CreatePost(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := user.(*domainUser.UserModel)

	// 1) まずmultipartをパース
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form", "details": err.Error()})
		return
	}

	// 2) テキスト系をDTOにマッピング:
	var input dto.CreatePostInput
	// 例: c.ShouldBind(&input) でも構わないが、multipart + JSON 混在の場合は工夫が必要
	// ここではシンプルに ctx.PostForm() して manual で入れるか
	input.Title = ctx.PostForm("title")
	input.Description = ctx.PostForm("description")
	// genres, skillsは PostFormArray() で取得して、inputへ
	input.Genres = ctx.PostFormArray("genres")
	input.Skills = ctx.PostFormArray("skills")

	// 3) 画像はmultipart.FileHeaderで受け取る
	form, _ := ctx.MultipartForm()
	fileHeaders := form.File["images"] // []*multipart.FileHeader

	// 4) サービスに「DTO + 画像ファイル群 + userID」を渡す
	err := c.portfolioService.CreatePost(input, fileHeaders, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

func (c *PortfolioController) GetPostsByUserID(ctx *gin.Context) {
	// ユーザー認証情報を取得
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := user.(*domainUser.UserModel)

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

func (c *PortfolioController) GetPostByID(ctx *gin.Context) {
	// URLパラメータ :id を取得
	idStr := ctx.Param("id")
	// 整数にパース
	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	postID := uint(idUint64)

	// サービスを呼び出して該当のPostを取得
	post, err := c.portfolioService.GetPostByID(postID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found", "details": err.Error()})
		return
	}

	// 取得したPostをレスポンスとして返す
	ctx.JSON(http.StatusOK, gin.H{"post": post})
}
