package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web3_task_blog/internal/utils"
)

// SecretKey is loaded from jwt_utils.go
func JWTMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// �d /login �
		if c.Request.URL.Path == "/static/login.html" || c.Request.URL.Path == "/login" {
			c.Next()
			return
		}
		tokenString, err := c.Cookie("token")
		_, errs := utils.ValidateToken(tokenString)
		if tokenString == "" || err != nil || errs != nil {
			c.Redirect(http.StatusFound, "/static/login.html")
			c.Abort()
			return
		}
		c.Next()
	}
}
