package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web3_task_blog/internal/handler/dto"
	"web3_task_blog/internal/repository"
	"web3_task_blog/internal/repository/entity"
	"web3_task_blog/internal/utils"
)

// Index 首页处理函数
func Index(c *gin.Context) {
	// 获取文章列表
	postRepo := repository.NewPostRepository()
	posts, err := postRepo.GetAll(100, 0)
	if err != nil {
		// 如果没有文章，返回空列表而不是错误
		posts = []entity.Post{}
	}

	// 转换为DTO格式
	var postListDTOs []dto.PostListDTO
	for _, post := range posts {
		// 获取作者信息
		userRepo := repository.NewUserRepository()
		user, _ := userRepo.FindByID(post.UserID)

		// 截取内容预览（取前100个字符）
		contentPreview := post.Content
		if len(contentPreview) > 100 {
			contentPreview = contentPreview[:100] + "..."
		}

		postListDTO := dto.PostListDTO{
			ID:            post.ID,
			Title:         post.Title,
			Content:       contentPreview,
			UserID:        post.UserID,
			Username:      user.Username,
			CommentStatus: post.CommentStatus,
			CommentCount:  len(post.Comments),
			CreatedAt:     post.CreatedAt,
		}
		postListDTOs = append(postListDTOs, postListDTO)
	}

	// 检查用户是否已登录
	token, err := c.Cookie("token")
	isAuthenticated := false
	var currentUser entity.User
	if err == nil && token != "" {
		claims, err := utils.ValidateToken(token)
		if err == nil {
			isAuthenticated = true
			// 获取当前用户信息
			userRepo := repository.NewUserRepository()
			user, _ := userRepo.FindByID(claims.UserID)
			currentUser = *user
		}
	}

	// 渲染首页
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":           "博客首页",
		"posts":           postListDTOs,
		"isAuthenticated": isAuthenticated,
		"currentUser":     currentUser,
	})
}
