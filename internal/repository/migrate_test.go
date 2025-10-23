package repository

import (
	"testing"
	"web3_task_blog/internal/repository/entity"
)

func TestAutoMigrate(t *testing.T) {
	// 获取测试数据库连接
	db, err := GetTestDB()
	if err != nil {
		t.Errorf("Failed to connect to test database: %v", err)
		return
	}

	// 先删除已存在的表
	err = db.Migrator().DropTable(&entity.User{}, &entity.Post{}, &entity.Comment{})
	if err != nil {
		t.Errorf("Failed to drop tables: %v", err)
		return
	}

	// 测试自动迁移
	err = AutoMigrate(db)
	if err != nil {
		t.Errorf("AutoMigrate failed: %v", err)
		return
	}

	// 检查表是否存在
	// 尝试创建一个用户来验证表是否存在
	user := &entity.User{
		Username: "testuser",
		Password: "testpassword",
		Email:    "test@example.com",
	}

	err = db.Create(user).Error
	if err != nil {
		t.Errorf("Failed to create test user: %v", err)
		return
	}

	// 清理测试数据
	db.Delete(user)
}