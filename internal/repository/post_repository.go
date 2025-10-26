package repository

import (
	"web3_task_blog/internal/repository/entity"
)

// PostRepository 文章仓库接口
type PostRepository interface {
	Create(post *entity.Post) error
	GetByID(id uint) (*entity.Post, error)
	GetAll(limit, offset int) ([]entity.Post, error)
	Update(post *entity.Post) error
	Delete(id uint) error
	GetByUserID(userID uint) ([]entity.Post, error)
}

// postRepository 文章仓库实现
type postRepository struct {
}

// NewPostRepository 创建文章仓库实例
func NewPostRepository() PostRepository {
	return &postRepository{}
}

// Create 创建文章
func (r *postRepository) Create(post *entity.Post) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Create(post).Error
}

// GetByID 根据ID获取文章
func (r *postRepository) GetByID(id uint) (*entity.Post, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var post entity.Post
	err = db.Preload("Comments").Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetAll 获取所有文章
func (r *postRepository) GetAll(limit, offset int) ([]entity.Post, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var posts []entity.Post
	err = db.Preload("User").Limit(limit).Offset(offset).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// Update 更新文章
func (r *postRepository) Update(post *entity.Post) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Save(post).Error
}

// Delete 删除文章
func (r *postRepository) Delete(id uint) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Delete(&entity.Post{}, id).Error
}

// GetByUserID 根据用户ID获取文章
func (r *postRepository) GetByUserID(userID uint) ([]entity.Post, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var posts []entity.Post
	err = db.Where("user_id = ?", userID).Preload("User").Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}