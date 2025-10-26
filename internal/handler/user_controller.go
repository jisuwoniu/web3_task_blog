package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"web3_task_blog/internal/handler/dto"
	"web3_task_blog/internal/repository"
	"web3_task_blog/internal/repository/entity"
	"web3_task_blog/internal/utils"
)

func Login(c *gin.Context) {
	var userDTO dto.UserDTO

	// 检查Content-Type，决定使用哪种方式获取数据
	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		// JSON格式请求
		var jsonReq struct {
			Username string `json:"user_name"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		userDTO.Username = jsonReq.Username
		userDTO.Password = jsonReq.Password
	} else {
		// 表单格式请求
		userDTO.Username = c.PostForm("user_name")
		userDTO.Password = c.PostForm("password")
	}

	if userDTO.Username == "" || userDTO.Password == "" {
		if contentType == "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		} else {
			c.Redirect(http.StatusFound, "/login.html?error=invalid_credentials")
		}
		return
	}

	repo := repository.NewUserRepository()
	userEntity, err := repo.FindByUsername(userDTO.Username)
	if err != nil {
		// 检查错误是否为表不存在
		if err.Error() == "Table 'users' doesn't exist" || err.Error() == "Error 1146: Table 'myapp.users' doesn't exist" {
			// 尝试自动建表
			if migrateErr := repository.AutoMigrate(); migrateErr != nil {
				if contentType == "application/json" {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database tables"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database tables"})
				}
				return
			}
		} else {
			if contentType == "application/json" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			} else {
				c.Redirect(http.StatusFound, "/login.html?error=login_failed")
			}
			return
		}
	}

	// 验证密码
	if !utils.CheckPasswordHash(userDTO.Password, userEntity.Password) {
		if contentType == "application/json" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.Redirect(http.StatusFound, "/login.html?error=login_failed")
		}
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(userEntity.ID)
	if err != nil {
		if contentType == "application/json" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		}
		return
	}
	c.SetCookie("token", token, 3600, "/", "", false, true)
	if contentType == "application/json" {
		// 返回JSON响应并设置Cookie
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"user": gin.H{
				"id":       userEntity.ID,
				"username": userEntity.Username,
				"email":    userEntity.Email,
			},
		})
	} else {
		// 设置Cookie并重定向
		c.Redirect(http.StatusFound, "/")
	}
}

// Register 用户注册处理函数
func Register(c *gin.Context) {
	var userDTO dto.UserDTO

	// 判断请求类型
	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		// 处理JSON请求
		if err := c.ShouldBindJSON(&userDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 处理表单请求
		userDTO.Username = c.PostForm("user_name")
		userDTO.Password = c.PostForm("password")
		userDTO.Email = c.PostForm("email")
	}

	// 验证必填字段
	if userDTO.Username == "" || userDTO.Password == "" || userDTO.Email == "" {
		if contentType == "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username, password and email are required"})
		} else {
			c.Redirect(http.StatusFound, "/static/register.html?error=missing_fields")
		}
		return
	}

	// 检查用户名是否已存在
	repo := repository.NewUserRepository()
	_, err := repo.FindByUsername(userDTO.Username)
	if err == nil {
		// 用户名已存在
		if contentType == "application/json" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.Redirect(http.StatusFound, "/static/register.html?error=username_exists")
		}
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		if contentType == "application/json" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		} else {
			c.Redirect(http.StatusFound, "/static/register.html?error=server_error")
		}
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
		if contentType == "application/json" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		} else {
			c.Redirect(http.StatusFound, "/static/register.html?error=server_error")
		}
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(userEntity.ID)
	if err != nil {
		if contentType == "application/json" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		} else {
			c.Redirect(http.StatusFound, "/static/register.html?error=server_error")
		}
		return
	}

	// 设置cookie
	c.SetCookie("token", token, 3600, "/", "", false, true)

	// 返回响应
	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"user_id": userEntity.ID,
		})
	} else {
		c.Redirect(http.StatusFound, "/static/hello.html")
	}
}

// ProfilePage 个人资料页面
func ProfilePage(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "User not authenticated"})
		return
	}

	// 将interface{}类型转换为uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case uint32:
		uid = uint(v)
	default:
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Invalid user ID type"})
		return
	}

	// 获取用户信息
	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindByID(uid)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "User not found"})
		return
	}

	// 获取用户的文章
	postRepo := repository.NewPostRepository()
	posts, err := postRepo.GetByUserID(uid)
	if err != nil {
		// 如果没有文章，返回空列表而不是错误
		posts = []entity.Post{}
	}

	// 转换为DTO格式
	var postListDTOs []dto.PostListDTO
	for _, post := range posts {
		// 截取内容预览（取前100个字符）
		contentPreview := post.Content
		if len(contentPreview) > 100 {
			contentPreview = contentPreview[:100] + "..."
		}

		postListDTO := dto.PostListDTO{
			ID:            post.ID,
			Title:         post.Title,
			Content:       contentPreview,
			UserID:        post.UserID,
			Username:      user.Username,
			CommentStatus: post.CommentStatus,
			CommentCount:  len(post.Comments),
			CreatedAt:     post.CreatedAt,
		}
		postListDTOs = append(postListDTOs, postListDTO)
	}

	// 渲染个人资料页面
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"title":      "个人资料",
		"user":       user,
		"posts":      postListDTOs,
		"postsCount": len(posts),
	})
}

// GetProfile 获取用户个人资料API
func GetProfile(c *gin.Context) {
	// 优先从Authorization header获取token
	tokenString := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// 检查是否是Bearer token
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:] // 去掉"Bearer "前缀
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}
	} else {
		// 如果header中没有token，从cookie获取
		var err error
		tokenString, err = c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			return
		}
	}

	// 解析token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// 获取用户信息
	repo := repository.NewUserRepository()
	user, err := repo.FindByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 转换为DTO，包含更多信息
	userDTO := dto.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userDTO,
	})
}

// GetUser 获取单个用户信息
func GetUser(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 获取用户信息
	repo := repository.NewUserRepository()
	user, err := repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 转换为DTO，包含更多信息
	userDTO := dto.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, userDTO)
}

// GetUserStatus 获取用户登录状态
func GetUserStatus(c *gin.Context) {
	// 尝试从cookie获取token
	tokenString, err := c.Cookie("token")

	// 如果cookie中没有token，尝试从Authorization header获取
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:] // 去掉"Bearer "前缀
		}
	}

	// 如果没有token，返回未登录状态
	if tokenString == "" {
		c.JSON(http.StatusOK, gin.H{
			"loggedIn": false,
		})
		return
	}

	// 验证token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		// token无效，返回未登录状态
		c.JSON(http.StatusOK, gin.H{
			"loggedIn": false,
		})
		return
	}

	// 获取用户信息
	repo := repository.NewUserRepository()
	user, err := repo.FindByID(claims.UserID)
	if err != nil {
		// 用户不存在，返回未登录状态
		c.JSON(http.StatusOK, gin.H{
			"loggedIn": false,
		})
		return
	}

	// 返回登录状态和用户信息
	c.JSON(http.StatusOK, gin.H{
		"loggedIn": true,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 清除cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
