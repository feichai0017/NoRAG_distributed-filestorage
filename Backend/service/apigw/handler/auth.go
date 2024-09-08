package handler

import (
	dblayer "cloud_distributed_storage/Backend/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

// Authorize 拦截器
func Authorize() gin.HandlerFunc {
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

		username, err := IsTokenValid(token)
		log.Printf("username: %s", username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将 username 添加到请求的上下文中
		c.Set("username", username)
		c.Next()

	}

}

func IsTokenValid(token string) (string, error) {
	if len(token) <= 8 {
		return "", fmt.Errorf("invalid token length")
	}

	// 验证 token 时效性
	tokenTimestamp := token[len(token)-8:]
	var ts int64
	_, err := fmt.Sscanf(tokenTimestamp, "%x", &ts)
	if err != nil {
		return "", fmt.Errorf("invalid token timestamp")
	}

	tokenTime := time.Unix(ts, 0)
	if time.Since(tokenTime) > 2*time.Hour {
		return "", fmt.Errorf("token expired")
	}
	username := dblayer.QueryUserByToken(token)
	if len(username) == 0 {
		return "", fmt.Errorf("username not found: %v", err)
	}
	return username, nil
}
