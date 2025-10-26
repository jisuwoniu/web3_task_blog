package dto

import (
	"time"
)

// PostDTO 文章数据传输对象
type PostDTO struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	UserID        uint      `json:"user_id"`
	Username      string    `json:"username"`       // 作者用户名
	CommentStatus string    `json:"comment_status"` // 评论状态
	CommentCount  int       `json:"comment_count"`  // 评论数量
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PostListDTO 文章列表数据传输对象
type PostListDTO struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"` // 截取的内容预览
	UserID        uint      `json:"user_id"`
	Username      string    `json:"username"`       // 作者用户名
	CommentStatus string    `json:"comment_status"` // 评论状态
	CommentCount  int       `json:"comment_count"`  // 评论数量
	CreatedAt     time.Time `json:"created_at"`
}

// PostCreateDTO 创建文章数据传输对象
type PostCreateDTO struct {
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	CommentStatus string `json:"comment_status"`
}

// PostUpdateDTO 更新文章数据传输对象
type PostUpdateDTO struct {
	ID            uint   `json:"id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	CommentStatus string `json:"comment_status"`
}

// PostDetailDTO 文章详情数据传输对象
type PostDetailDTO struct {
	ID            uint         `json:"id"`
	Title         string       `json:"title"`
	Content       string       `json:"content"`
	UserID        uint         `json:"user_id"`
	Username      string       `json:"username"`       // 作者用户名
	CommentStatus string       `json:"comment_status"` // 评论状态
	Comments      []CommentDTO `json:"comments"`       // 评论列表
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
