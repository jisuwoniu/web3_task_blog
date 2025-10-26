package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web3_task_blog/internal/handler/dto"
	"web3_task_blog/internal/repository"
	"web3_task_blog/internal/repository/entity"
	"web3_task_blog/internal/utils"
)

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
