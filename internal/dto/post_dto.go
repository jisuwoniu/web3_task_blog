package dto

// PostDTO 文章数据传输对象
type PostDTO struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Category string `json:"category"`
}

// CommentDTO 评论数据传输对象
type CommentDTO struct {
	Content string `json:"content" binding:"required"`
}