package repository

import (
	"web3_task_blog/internal/repository/entity"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.AutoMigrate(
		&entity.User{},
		&entity.Post{},
		&entity.Comment{},
	)
}