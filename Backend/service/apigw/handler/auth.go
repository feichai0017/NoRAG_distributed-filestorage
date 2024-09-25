package handler

import (
	rPool "cloud_distributed_storage/Backend/cache/redis"
	"fmt"
	"github.com/garyburd/redigo/redis"
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

	// 从 Redis 中获取用户名
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	username, err := redis.String(rConn.Do("GET", fmt.Sprintf("session*%s", token)))
	if err != nil {
		if err == redis.ErrNil {
			return "", fmt.Errorf("token not found or expired")
		}
		return "", fmt.Errorf("error querying Redis: %v", err)
	}

	// 验证 token 时效性 (可选，因为 Redis 已经处理了过期)
	tokenTimestamp := token[len(token)-8:]
	var ts int64
	_, err = fmt.Sscanf(tokenTimestamp, "%x", &ts)
	if err != nil {
		return "", fmt.Errorf("invalid token timestamp")
	}

	tokenTime := time.Unix(ts, 0)
	if time.Since(tokenTime) > 24*time.Hour {
		// 如果 token 已过期，从 Redis 中删除它
		_, err = rConn.Do("DEL", fmt.Sprintf("session*%s", token))
		if err != nil {
			fmt.Printf("Error deleting expired token from Redis: %v\n", err)
		}
		return "", fmt.Errorf("token expired")
	}

	return username, nil
}
