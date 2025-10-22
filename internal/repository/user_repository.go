package repository

import (
	"gorm.io/gorm"
	"web3_task_blog/internal/repository/entity"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint32) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteUser(id uint32) error {
	return r.db.Delete(&entity.User{}, id).Error
}
