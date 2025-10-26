package repository

import (
	"web3_task_blog/internal/repository/entity"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var user entity.User
	err = db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*entity.User, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	user := &entity.User{}
	err = db.First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id uint) (*entity.User, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	user := &entity.User{}
	err = db.First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *entity.User) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Save(user).Error
}

func (r *UserRepository) DeleteUser(id uint32) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	return db.Delete(&entity.User{}, id).Error
}