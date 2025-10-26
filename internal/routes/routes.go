package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web3_task_blog/internal/handler"
	"web3_task_blog/internal/middleware"
)

// SetupRoutes 配置路由
func SetupRoutes() *gin.Engine {
	r := gin.Default()
	// 配置静态文件服务
	r.Static("/static", "./web/static")
	//
	//// 加载HTML模板
	r.LoadHTMLGlob("web/static/*.html")
	// 公开路由
	public := r.Group("/")
	{
		// 页面路由
		public.GET("/", handler.Index)
		public.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		})
		public.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", gin.H{})
		})
		public.GET("/post/:id", handler.PostDetailPage)

		// API路由
		public.POST("/login", handler.Login)
		public.POST("/api/login", handler.Login)
		public.POST("/register", handler.Register)
		public.POST("/api/register", handler.Register)
		public.GET("/api/posts", handler.GetPosts)
		public.GET("/api/post/:id", handler.GetPost)
		public.GET("/api/post/:id/comments", handler.GetComments)
		public.GET("/api/user/:id/posts", handler.GetPostsByUserID)
		public.GET("/api/user/status", handler.GetUserStatus)
	}

	// 需要认证的路由
	authorized := r.Group("/")
	authorized.Use(middleware.JWTMiddle())
	{
		// 页面路由
		authorized.GET("/profile", handler.ProfilePage)
		authorized.GET("/posts/create-post", handler.CreatePostPage)
		authorized.GET("/posts/edit-post/:id", handler.EditPostPage)
		// API路由
		authorized.GET("/api/user/:id", handler.GetUser)
		authorized.POST("/api/posts", handler.CreatePost)
		authorized.PUT("/api/post/:id/update", handler.UpdatePost)
		authorized.DELETE("/api/post/:id/delete", handler.DeletePost)
		authorized.POST("/api/post/:id/comments/add", handler.CreateComment)
		authorized.DELETE("/api/comments/:id", handler.DeleteComment)
		authorized.GET("/api/profile", handler.GetProfile)
	}

	return r
}
