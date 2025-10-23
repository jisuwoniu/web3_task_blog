package repository

import (
	"gorm.io/gorm"
	"web3_task_blog/internal/repository/entity"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Post{},
		&entity.Comment{},
	)
}