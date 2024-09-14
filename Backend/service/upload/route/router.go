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

	// 文件上传相关接口
	r.POST("/file/upload", api.UploadHandler)
	// 秒传接口
	r.POST("/file/fastupload", api.TryFastUploadHandler)

	// 分块上传接口
	r.POST("/file/mpupload/init", api.InitialMultipartUploadHandler)
	r.POST("/file/mpupload/uppart", api.UploadPartHandler)
	r.POST("/file/mpupload/complete", api.CompleteUploadHandler)
	r.POST("file/mpupload/cancel", api.CancelUploadHandler)
	r.POST("file/mpupload/status", api.MultipartUploadStatusHandler)
	r.POST("file/mpupload/multi", api.MultiDownloadHandler)

	return r
}
