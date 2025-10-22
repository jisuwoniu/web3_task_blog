package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web3_task_blog/internal/handler/dto"
	"web3_task_blog/internal/repository"
	"web3_task_blog/internal/utils"
)

func Login(c *gin.Context) {
	var userDTO dto.UserDTO

	userDTO.Username = c.PostForm("user_name")
	userDTO.Password = c.PostForm("password")
	if userDTO.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	db, err := repository.GetDB()
	repo := repository.NewUserRepository(db)
	userEntity, err := repo.FindByUsername(userDTO.Username)
	if err != nil || userEntity.Password != userDTO.Password {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	token, err := utils.GenerateToken(userEntity.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetCookie("token", token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
