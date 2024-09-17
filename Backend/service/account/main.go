package main

import (
	"cloud_distributed_storage/Backend/common"
	"cloud_distributed_storage/Backend/service/account/handler"
	userProto "cloud_distributed_storage/Backend/service/account/proto"
	dbproxy "cloud_distributed_storage/Backend/service/dbproxy/client"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"log"
	"time"
)

func main() {
	// 创建 Consul 注册中心
	reg := consul.NewRegistry(registry.Addrs("localhost:8500"))
	// 创建服务
	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

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
