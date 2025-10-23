package entity

import "gorm.io/gorm"

type Comment struct {
	CommentInfo string `gorm:"comment:评论内容"`
	PostID      uint   `gorm:"comment:文章编号"`
	gorm.Model
}
