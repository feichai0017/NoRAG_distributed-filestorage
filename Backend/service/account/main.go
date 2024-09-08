package main

import (
	"cloud_distributed_storage/Backend/common"
	"cloud_distributed_storage/Backend/service/account/handler"
	userProto "cloud_distributed_storage/Backend/service/account/proto"
	registry "cloud_distributed_storage/Backend/service/registry"
	"github.com/asim/go-micro/v3"
	"log"
	"time"
)

func main() {
	// Create Consul registry
	consulReg := registry.GetConsulRegistry()

	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Registry(consulReg),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	service.Init()

	// 注册处理程序
	if err := userProto.RegisterUserServiceHandler(service.Server(), new(handler.User)); err != nil {
		log.Fatalf("Failed to register handler: %v", err)
	}

	// 运行服务
	log.Println("Starting user service")
	if err := service.Run(); err != nil {
		log.Fatalf("Service failed to run: %v", err)
	}
}
