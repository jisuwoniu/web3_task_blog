package dto

import "time"

// CommentDTO 评论数据传输对象
type CommentDTO struct {
	ID          uint      `json:"id"`
	CommentInfo string    `json:"comment_info"`
	PostID      uint      `json:"post_id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"` // 评论者用户名
	CreatedAt   time.Time `json:"created_at"`
}

// CommentCreateDTO 创建评论数据传输对象
type CommentCreateDTO struct {
	CommentInfo string `json:"comment_info" binding:"required"`
	PostID      uint   `json:"post_id" binding:"required"`
}
