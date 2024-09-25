package route

import (
	"cloud_distributed_storage/Backend/service/apigw/handler"
	"cloud_distributed_storage/Backend/service/apigw/middleware"
	"github.com/gin-gonic/gin"
)

// Router: gateway api router
func Router() *gin.Engine {
	router := gin.Default()

	// 使用CORS中间件
	router.Use(middleware.CORSMiddleware())

	router.POST("/user/signup", handler.SignupHandler)

	router.POST("/user/login", handler.SignInHandler)

	// 需要认证的路由
	auth := router.Group("/")
	auth.Use(handler.Authorize())
	{
		auth.GET("/user/info", handler.UserInfoHandler)
		auth.GET("/user/logout", handler.SignOutHandler)
		auth.POST("/user/delete", handler.DeleteUserHandler)
		auth.POST("/file/query", handler.FileQueryHandler)
		auth.POST("/file/update", handler.FileMetaUpdateHandler)
	}

	return router
}
