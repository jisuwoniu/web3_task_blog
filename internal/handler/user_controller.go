package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web3_task_blog/internal/handler/dto"
	"web3_task_blog/internal/repository"
	"web3_task_blog/internal/repository/entity"
	"web3_task_blog/internal/utils"
)

func Login(c *gin.Context) {
	var userDTO dto.UserDTO

	userDTO.Username = c.PostForm("user_name")
	userDTO.Password = c.PostForm("password")
	if userDTO.Username == "" || userDTO.Password == "" {
		c.Redirect(http.StatusFound, "/static/login.html?error=invalid_credentials")
		return
	}

	db, err := repository.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	repo := repository.NewUserRepository(db)
	userEntity, err := repo.FindByUsername(userDTO.Username)
	if err != nil {
		// 检查错误是否为表不存在
		if err.Error() == "Table 'users' doesn't exist" || err.Error() == "Error 1146: Table 'myapp.users' doesn't exist" {
			// 尝试自动建表
			if migrateErr := repository.AutoMigrate(db); migrateErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database tables"})
				return
			}
		} else {
			c.Redirect(http.StatusFound, "/static/login.html?error=login_failed")
			return
		}
	}

	// 验证密码
	if !utils.CheckPasswordHash(userDTO.Password, userEntity.Password) {
		c.Redirect(http.StatusFound, "/static/login.html?error=login_failed")
		return
	}

	token, err := utils.GenerateToken(uint32(userEntity.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetCookie("token", token, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/static/hello.html")
}

func Register(c *gin.Context) {
	var userDTO dto.UserDTO

	userDTO.Username = c.PostForm("user_name")
	userDTO.Password = c.PostForm("password")
	userDTO.Email = c.PostForm("email")

	// 验证必填字段
	if userDTO.Username == "" || userDTO.Password == "" || userDTO.Email == "" {
		c.Redirect(http.StatusFound, "/static/register.html?error=missing_fields")
		return
	}

	db, err := repository.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	repo := repository.NewUserRepository(db)

	// 检查用户名是否已存在
	_, err = repo.FindByUsername(userDTO.Username)
	if err != nil {
		// 检查错误是否为表不存在
		if err.Error() == "Table 'users' doesn't exist" || err.Error() == "Error 1146: Table 'myapp.users' doesn't exist" {
			// 尝试自动建表
			if migrateErr := repository.AutoMigrate(db); migrateErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database tables"})
				return
			}
			// 再次尝试查询用户
			_, err = repo.FindByUsername(userDTO.Username)
			if err == nil {
				// 如果查询成功，说明用户已存在
				c.Redirect(http.StatusFound, "/static/register.html?error=username_exists")
				return
			}
		} else if err.Error() != "record not found" {
			// 如果是其他错误，直接返回
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		// 用户名已存在
		c.Redirect(http.StatusFound, "/static/register.html?error=username_exists")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建用户实体
	userEntity := &entity.User{
		Username: userDTO.Username,
		Password: hashedPassword,
		Email:    userDTO.Email,
	}

	// 保存用户
	err = repo.CreateUser(userEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 注册成功，重定向到登录页面
	c.Redirect(http.StatusFound, "/static/login.html?success=registration_successful")
}