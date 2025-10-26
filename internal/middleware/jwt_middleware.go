package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"web3_task_blog/internal/utils"
)

func JWTMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从Authorization header获取token
		tokenString := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// 检查是否是Bearer token
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = authHeader[7:] // 去掉"Bearer "前缀
			}
		}

		// 如果header中没有token，从cookie获取
		if tokenString == "" {
			var err error
			tokenString, err = c.Cookie("token")
			if err != nil {
				// 如果是API请求，返回JSON错误
				if strings.HasPrefix(c.Request.URL.Path, "/api/") {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
					c.Abort()
					return
				}
				// 如果是页面请求，重定向到登录页面
				c.Redirect(http.StatusFound, "/login?error=login_required")
				c.Abort()
				return
			}
		}

		// 验证token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			// 如果是API请求，返回JSON错误
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
			// 如果是页面请求，重定向到登录页面
			c.Redirect(http.StatusFound, "/login?error=invalid_token")
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
