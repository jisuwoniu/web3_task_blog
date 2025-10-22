package routes

import (
	"github.com/gin-gonic/gin"
	//"net/http"
	"web3_task_blog/internal/handler"
	"web3_task_blog/internal/middleware"
	//"web3_task_blog/internal/repository"
	//"web3_task_blog/internal/utils"
)

// SetupRoutes 配置路由
func SetupRoutes() *gin.Engine {
	r := gin.Default()
	// 需要认证的路由
	r.Use(middleware.JWTMiddle())
	r.Static("/static", "./web/static")
	//authorized := r.Group("/api")
	//authorized.Use(middleware.JWTMiddle())
	//{
	//	authorized.GET("/profile", func(c *gin.Context) {
	//		userID := c.GetInt("user_id")
	//		c.JSON(http.StatusOK, gin.H{"user_id": userID, "message": "Welcome to your profile!"})
	//	})
	//}

	// 公开路由
	r.POST("/login", handler.Login)
	return r
}
