package route

import (
	"cloud_distributed_storage/Backend/service/upload/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 使用gin插件支持跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://192.168.0.200:3001"}, // []string{"http://localhost:8081"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
		// AllowCredentials: true,
	}))

	r.POST("/file/upload", api.UploadHandler)

	return r
}
