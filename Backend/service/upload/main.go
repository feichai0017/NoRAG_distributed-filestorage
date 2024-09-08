package main

import (
	"cloud_distributed_storage/Backend/common"
	"cloud_distributed_storage/Backend/config"
	"cloud_distributed_storage/Backend/mq"
	"cloud_distributed_storage/Backend/route"
	dbproxy "cloud_distributed_storage/Backend/service/dbproxy/client"
	registry "cloud_distributed_storage/Backend/service/registry"
	cfg "cloud_distributed_storage/Backend/service/upload/config"
	upProto "cloud_distributed_storage/Backend/service/upload/proto"
	upRpc "cloud_distributed_storage/Backend/service/upload/rpc"
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func startRPCService() {
	// Create Consul registry
	consulReg := registry.GetConsulRegistry()

	service := micro.NewService(
		micro.Name("go.micro.service.upload"), // 服务名称
		micro.Registry(consulReg),
		micro.RegisterTTL(time.Second*10),     // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) error {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
			return nil
		}),
	)

	// 初始化dbproxy client
	dbproxy.Init(service)
	// 初始化mq client
	mq.Init()

	err := upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err != nil {
		log.Fatalf("Failed to register handler: %v", err)
	}
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	err := router.Run(cfg.UploadServiceHost)
	if err != nil {
		log.Fatal(err)
	}
	// service := web.NewService(
	// 	web.Name("go.micro.web.upload"),
	// 	web.Handler(router),
	// 	web.RegisterTTL(10*time.Second),
	// 	web.RegisterInterval(5*time.Second),
	// )
	// if err := service.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := service.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {
	err := os.MkdirAll(config.TempLocalRootDir, 0777)
	if err != nil {
		log.Println(err)
		return
	}

	err = os.MkdirAll(config.TempPartRootDir, 0777)
	if err != nil {
		log.Println(err)
		return
	}

	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
