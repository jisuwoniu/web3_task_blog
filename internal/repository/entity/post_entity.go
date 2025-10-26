package entity

import "gorm.io/gorm"

type Post struct {
	Title         string    `gorm:"comment:文章标题;type:varchar(128)"`
	Content       string    `gorm:"comment:文章内容"`
	Category      string    `gorm:"comment:文章分类;type:varchar(50)"`
	UserID        uint      `gorm:"comment:文章作者userId"`
	CommentStatus string    `gorm:"column:comment_status;comment:评论状态"` // 评论状态（如 "有评论" 或 "无评论"）
	User          User      `gorm:"foreignKey:UserID"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
	gorm.Model
}
