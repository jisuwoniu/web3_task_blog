package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"web3_task_blog/internal/handler/dto"
	"web3_task_blog/internal/repository"
	"web3_task_blog/internal/repository/entity"
	"web3_task_blog/internal/utils"
)

// GetPosts 获取文章列表
func GetPosts(c *gin.Context) {
	postRepo := repository.NewPostRepository()
	posts, err := postRepo.GetAll(100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
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

	c.JSON(http.StatusOK, gin.H{"posts": postListDTOs})
}

// GetPost 获取单篇文章详情
func GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	postRepo := repository.NewPostRepository()
	post, err := postRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 获取作者信息
	userRepo := repository.NewUserRepository()
	user, _ := userRepo.FindByID(post.UserID)

	// 转换评论为DTO格式
	var commentDTOs []dto.CommentDTO
	for _, comment := range post.Comments {
		commentUser, _ := userRepo.FindByID(comment.UserID)
		commentDTO := dto.CommentDTO{
			ID:          comment.ID,
			CommentInfo: comment.CommentInfo,
			PostID:      comment.PostID,
			UserID:      comment.UserID,
			Username:    commentUser.Username,
			CreatedAt:   comment.CreatedAt,
		}
		commentDTOs = append(commentDTOs, commentDTO)
	}

	// 转换为DTO格式
	postDetailDTO := dto.PostDetailDTO{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		UserID:        post.UserID,
		Username:      user.Username,
		CommentStatus: post.CommentStatus,
		Comments:      commentDTOs,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}

	c.JSON(http.StatusOK, postDetailDTO)
}

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var postCreateDTO dto.PostCreateDTO
	if err := c.ShouldBindJSON(&postCreateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT中获取用户ID
	id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 将interface{}类型转换为uint
	var uid uint
	switch v := id.(type) {
	case uint:
		uid = v
	case uint32:
		uid = uint(v)
	case *utils.Claims:
		uid = uint(v.UserID)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// 创建文章实体
	postEntity := &entity.Post{
		Title:         postCreateDTO.Title,
		Content:       postCreateDTO.Content,
		UserID:        uid,
		CommentStatus: postCreateDTO.CommentStatus,
	}

	postRepo := repository.NewPostRepository()
	err := postRepo.Create(postEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post_id": postEntity.ID})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {

	var postUpdateDTO dto.PostUpdateDTO
	if err := c.ShouldBindJSON(&postUpdateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 从JWT中获取用户ID
	id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 将interface{}类型转换为uint
	var uid uint
	switch v := id.(type) {
	case uint:
		uid = v
	case uint32:
		uid = uint(v)
	case *utils.Claims:
		uid = uint(v.UserID)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	postRepo := repository.NewPostRepository()
	// 检查文章是否存在
	post, err := postRepo.GetByID(postUpdateDTO.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 检查用户是否有权限修改文章
	if post.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this post"})
		return
	}

	// 更新文章
	post.Title = postUpdateDTO.Title
	post.Content = postUpdateDTO.Content
	post.CommentStatus = postUpdateDTO.CommentStatus

	err = postRepo.Update(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 将interface{}类型转换为uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case uint32:
		uid = uint(v)
	case *utils.Claims:
		uid = uint(v.UserID)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	postRepo := repository.NewPostRepository()
	// 检查文章是否存在
	post, err := postRepo.GetByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 检查用户是否有权限删除文章
	if post.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this post"})
		return
	}

	err = postRepo.Delete(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// 从路径参数中获取文章ID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 将interface{}类型转换为uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case uint32:
		uid = uint(v)
	case *utils.Claims:
		uid = uint(v.UserID)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postRepo := repository.NewPostRepository()
	// 检查文章是否存在
	_, err = postRepo.GetByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 创建评论实体
	commentEntity := &entity.Comment{
		CommentInfo: req.Content,
		PostID:      uint(postID),
		UserID:      uid,
	}

	commentRepo := repository.NewCommentRepository()
	err = commentRepo.Create(commentEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 更新文章的评论状态
	post, _ := postRepo.GetByID(uint(postID))
	post.CommentStatus = "有评论"
	postRepo.Update(post)

	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully", "comment_id": commentEntity.ID})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 将interface{}类型转换为uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case uint32:
		uid = uint(v)
	case *utils.Claims:
		uid = uint(v.UserID)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	commentRepo := repository.NewCommentRepository()
	// 检查评论是否存在
	comment, err := commentRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// 检查用户是否有权限删除评论
	if comment.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
		return
	}

	err = commentRepo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// 检查文章是否还有其他评论
	postRepo := repository.NewPostRepository()
	post, _ := postRepo.GetByID(comment.PostID)
	comments, _ := commentRepo.FindByPostID(comment.PostID)
	if len(comments) == 0 {
		post.CommentStatus = "无评论"
		postRepo.Update(post)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// GetPostsByUserID 获取指定用户的文章列表
func GetPostsByUserID(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	postRepo := repository.NewPostRepository()
	posts, err := postRepo.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
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

	c.JSON(http.StatusOK, gin.H{"posts": postListDTOs})
}

// PostDetailPage 文章详情页面
func PostDetailPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid post ID"})
		return
	}

	postRepo := repository.NewPostRepository()
	post, err := postRepo.GetByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}

	// 获取作者信息
	userRepo := repository.NewUserRepository()
	user, _ := userRepo.FindByID(post.UserID)

	// 转换评论为DTO格式
	var commentDTOs []dto.CommentDTO
	for _, comment := range post.Comments {
		commentUser, _ := userRepo.FindByID(comment.UserID)
		commentDTO := dto.CommentDTO{
			ID:          comment.ID,
			CommentInfo: comment.CommentInfo,
			PostID:      comment.PostID,
			UserID:      comment.UserID,
			Username:    commentUser.Username,
			CreatedAt:   comment.CreatedAt,
		}
		commentDTOs = append(commentDTOs, commentDTO)
	}

	// 转换为DTO格式
	postDetailDTO := dto.PostDetailDTO{
		ID:            post.ID,
		Title:         post.Title,
		Content:       strings.ReplaceAll(post.Content, "\n", "<br>"),
		UserID:        post.UserID,
		Username:      user.Username,
		CommentStatus: post.CommentStatus,
		Comments:      commentDTOs,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}

	// 检查用户是否已登录
	token, err := c.Cookie("token")
	isAuthenticated := false
	if err == nil && token != "" {
		_, err = utils.ValidateToken(token)
		if err == nil {
			isAuthenticated = true
		}
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"post":            postDetailDTO,
		"isAuthenticated": isAuthenticated,
	})
}

// GetComments 获取文章评论
func GetComments(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// 获取文章
	postRepo := repository.NewPostRepository()
	post, err := postRepo.GetByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 获取评论
	commentRepo := repository.NewCommentRepository()
	comments, err := commentRepo.FindByPostID(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	// 获取评论用户信息
	userRepo := repository.NewUserRepository()
	var commentDTOs []dto.CommentDTO
	for _, comment := range comments {
		user, _ := userRepo.FindByID(comment.UserID)
		commentDTO := dto.CommentDTO{
			ID:          comment.ID,
			CommentInfo: comment.CommentInfo,
			PostID:      comment.PostID,
			UserID:      comment.UserID,
			Username:    user.Username,
			CreatedAt:   comment.CreatedAt,
		}
		commentDTOs = append(commentDTOs, commentDTO)
	}

	// 返回评论列表
	c.JSON(http.StatusOK, gin.H{
		"post_id":    post.ID,
		"post_title": post.Title,
		"comments":   commentDTOs,
	})
}

// CreatePostPage 创建文章页面
func CreatePostPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create-post.html", gin.H{})
}

// EditPostPage 编辑文章页面
func EditPostPage(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid post ID"})
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "User not authenticated"})
		return
	}

	// 将interface{}类型转换为uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case uint32:
		uid = uint(v)
	case *utils.Claims:
		uid = uint(v.UserID)
	default:
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Invalid user ID type"})
		return
	}

	postRepo := repository.NewPostRepository()
	post, err := postRepo.GetByID(uint(postID))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}

	// 检查用户是否有权限编辑文章
	if post.UserID != uid {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "You don't have permission to edit this post"})
		return
	}

	c.HTML(http.StatusOK, "edit_post.html", gin.H{
		"post": post,
	})
}

// PostDetail 文章详情页面
func PostDetail(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid post ID"})
		return
	}

	postRepo := repository.NewPostRepository()
	post, err := postRepo.GetByID(uint(postID))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}

	// 获取作者信息
	userRepo := repository.NewUserRepository()
	author, err := userRepo.FindByID(post.UserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to get author info"})
		return
	}

	// 获取评论
	commentRepo := repository.NewCommentRepository()
	comments, err := commentRepo.FindByPostID(uint(postID))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to get comments"})
		return
	}

	// 获取评论用户信息
	var commentUsers []*entity.User
	for _, comment := range comments {
		user, err := userRepo.FindByID(comment.UserID)
		if err != nil {
			continue
		}
		commentUsers = append(commentUsers, user)
	}

	// 检查用户是否已登录
	userID, exists := c.Get("user_id")
	var isAuthor bool
	if exists {
		// 将interface{}类型转换为uint
		var uid uint
		switch v := userID.(type) {
		case uint:
			uid = v
		case uint32:
			uid = uint(v)
		case *utils.Claims:
			uid = uint(v.UserID)
		default:
			// 如果类型不匹配，设置isAuthor为false
			isAuthor = false
		}
		isAuthor = (post.UserID == uid)
	}

	c.HTML(http.StatusOK, "post_detail.html", gin.H{
		"post":         post,
		"author":       author,
		"comments":     comments,
		"commentUsers": commentUsers,
		"isAuthor":     isAuthor,
	})
}