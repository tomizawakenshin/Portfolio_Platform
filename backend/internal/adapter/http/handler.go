package http

import (
	"fmt"
	"net/http"
	"strconv"

	domainUser "backend/domain/user"
	"backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// Handler は PortfolioUsecase を受け取って HTTP リクエストを処理します
type Handler struct {
	PortfolioUsecase usecase.PortfolioUsecase
}

// NewHandler は Handler を生成します
func NewHandler(portfolioUsecase usecase.PortfolioUsecase) *Handler {
	return &Handler{PortfolioUsecase: portfolioUsecase}
}

// CreatePost は投稿を作成します
func (handler *Handler) CreatePost(c *gin.Context) {
	var createPostRequest CreatePostRequest
	// multipart/form-data のバインド
	if err := c.ShouldBind(&createPostRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 認証ユーザーを取得
	userRaw, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := userRaw.(*domainUser.UserModel)

	// DTO に変換
	createPostInput, fileHeaders := createPostRequest.ToDTO(currentUser.ID)

	// ユースケース呼び出し
	if err := handler.PortfolioUsecase.CreatePost(createPostInput, fileHeaders, currentUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

// GetPostByID は ID で単一投稿を取得します
func (handler *Handler) GetPostByID(c *gin.Context) {
	idString := c.Param("id")
	postID, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := handler.PortfolioUsecase.GetPostByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found", "details": err.Error()})
		return
	}

	// レスポンスDTOへマッピングして返却
	postResponse := NewPostResponse(post)
	c.JSON(http.StatusOK, postResponse)
}

// GetPostsByUserID はログインユーザーの投稿一覧を取得します
func (handler *Handler) GetPostsByUserID(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	currentUser := userRaw.(*domainUser.UserModel)

	posts, err := handler.PortfolioUsecase.GetPostsByUserID(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	postsResponse := NewPostsResponse(posts)
	c.JSON(http.StatusOK, postsResponse)
}

// GetAllPosts は全投稿を取得します
func (handler *Handler) GetAllPosts(c *gin.Context) {
	posts, err := handler.PortfolioUsecase.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	postsResponse := NewPostsResponse(posts)

	// ここで0番目を出力
	if len(postsResponse) > 0 {
		fmt.Printf("[DEBUG] postsResponse[0]: %+v\n", postsResponse[0])
	}

	c.JSON(http.StatusOK, postsResponse)
}
