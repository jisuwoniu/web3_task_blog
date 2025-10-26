package repository

import (
	"web3_task_blog/internal/repository/entity"
)

// CommentRepository 评论仓库接口
type CommentRepository interface {
	Create(comment *entity.Comment) error
	FindByID(id uint) (*entity.Comment, error)
	FindByPostID(postID uint) ([]entity.Comment, error)
	FindByUserID(userID uint) ([]entity.Comment, error)
	Update(comment *entity.Comment) error
	Delete(id uint) error
	FindAll() ([]entity.Comment, error)
}

// commentRepository 评论仓库实现
type commentRepository struct {
}

// NewCommentRepository 创建评论仓库实例
func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

// Create 创建评论
func (r *commentRepository) Create(comment *entity.Comment) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Create(comment).Error
}

// FindByID 根据ID查找评论
func (r *commentRepository) FindByID(id uint) (*entity.Comment, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var comment entity.Comment
	err = db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindByPostID 根据文章ID查找评论
func (r *commentRepository) FindByPostID(postID uint) ([]entity.Comment, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var comments []entity.Comment
	err = db.Where("post_id = ?", postID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// FindByUserID 根据用户ID查找评论
func (r *commentRepository) FindByUserID(userID uint) ([]entity.Comment, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var comments []entity.Comment
	err = db.Where("user_id = ?", userID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Update 更新评论
func (r *commentRepository) Update(comment *entity.Comment) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Save(comment).Error
}

// Delete 删除评论
func (r *commentRepository) Delete(id uint) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Delete(&entity.Comment{}, id).Error
}

// FindAll 查找所有评论
func (r *commentRepository) FindAll() ([]entity.Comment, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var comments []entity.Comment
	err = db.Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}