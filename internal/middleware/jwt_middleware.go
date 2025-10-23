package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web3_task_blog/internal/utils"
)

func JWTMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许访问登录页面、注册页面和静态资源
		if c.Request.URL.Path == "/static/login.html" ||
			c.Request.URL.Path == "/static/register.html" ||
			c.Request.URL.Path == "/login" ||
			c.Request.URL.Path == "/register" {
			c.Next()
			return
		}

		tokenString, err := c.Cookie("token")
		_, errs := utils.ValidateToken(tokenString)
		if tokenString == "" || err != nil || errs != nil {
			c.Redirect(http.StatusFound, "/static/login.html?error=login_failed")
			c.Abort()
			return
		}
		c.Next()
	}
}
