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

	router.GET("/user/info", handler.UserInfoHandler)

	return router
}
