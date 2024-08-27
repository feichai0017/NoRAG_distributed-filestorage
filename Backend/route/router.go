package route

import (
	"cloud_distributed_storage/Backend/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/api/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	r.POST("/api/file/upload/success", handler.HTTPInterceptor(handler.UploadSucHandler))
	r.POST("/api/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	r.POST("/api/file/meta/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	r.GET("/api/file/download", handler.HTTPInterceptor(handler.DownloadHandler))
	r.POST("/api/file/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	r.POST("/api/file/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))
	r.POST("/api/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	// redis upload api
	r.POST("/api/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	r.POST("/api/mpupload/uploadpart", handler.HTTPInterceptor(handler.UploadPartHandler))
	r.POST("/api/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))
	r.POST("api/mpupload/cancel", handler.HTTPInterceptor(handler.CancelUploadHandler))
	r.POST("api/mpupload/status", handler.HTTPInterceptor(handler.MultipartUploadStatusHandler))

	// user api
	r.POST("/api/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	r.POST("/api/user/login", handler.SignInHandler)
	r.POST("/api/user/signup", handler.SignupHandler)

	return r
}
