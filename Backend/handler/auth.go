package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// HTTPInterceptor 拦截器
func HTTPInterceptor(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("Received Authorization header: %s", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
			c.Abort()
			return
		}

		// 解析 Bearer token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("token: %s", token)

		username, err := ValidateTokenAndGetUsername(token)
		log.Printf("username: %s", username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将 username 添加到请求的上下文中
		c.Set("username", username)
		h(c)

	}

}
