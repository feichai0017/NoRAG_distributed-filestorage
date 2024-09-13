package main

import (
	"cloud_distributed_storage/Backend/common"
	dbproxy "cloud_distributed_storage/Backend/service/dbproxy/client"
	cfg "cloud_distributed_storage/Backend/service/download/config"
	dlProto "cloud_distributed_storage/Backend/service/download/proto"
	"cloud_distributed_storage/Backend/service/download/route"
	dlRpc "cloud_distributed_storage/Backend/service/download/rpc"
	"cloud_distributed_storage/Backend/service/registry"
	"fmt"
	"github.com/asim/go-micro/v3"
	"time"
)

func startRPCService() {
	consulReg := registry.GetConsulRegistry()

	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.Registry(consulReg),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	err := dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
