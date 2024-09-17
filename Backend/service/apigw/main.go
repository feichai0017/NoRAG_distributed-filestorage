package main

import (
	"cloud_distributed_storage/Backend/service/apigw/handler"
	"cloud_distributed_storage/Backend/service/apigw/route"
	"log"
	"time"
)

func main() {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err := handler.InitService(); err == nil {
			break
		}
		log.Printf("Failed to initialize service, retrying in 5 seconds... (Attempt %d/%d)", i+1, maxRetries)
		time.Sleep(5 * time.Second)
	}

	// 初始化路由
	r := route.Router()

	// 启动 HTTP 服务器
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
