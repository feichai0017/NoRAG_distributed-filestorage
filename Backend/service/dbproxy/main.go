package main

import (
	"cloud_distributed_storage/Backend/common"
	"cloud_distributed_storage/Backend/service/dbproxy/config"
	dbConn "cloud_distributed_storage/Backend/service/dbproxy/conn"
	dbProxy "cloud_distributed_storage/Backend/service/dbproxy/proto"
	dbRpc "cloud_distributed_storage/Backend/service/dbproxy/rpc"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/urfave/cli/v2"
	"log"
	"time"
)

func startRpcService() {
	// 创建 Consul 注册中心
	reg := consul.NewRegistry(registry.Addrs("localhost:8500"))
	//start rpc server
	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	service.Init(
		micro.Action(func(c *cli.Context) error {
			// 初始化数据库连接
			dbhost := c.String("dbhost")
			if len(dbhost) > 0 {
				log.Println("custom database host: ", dbhost)
				config.UpdateDBHost(dbhost)
			}
			return nil
		}),
	)
	// Init db connection
	dbConn.InitDBConn()

	// Register handler
	err := dbProxy.RegisterDBProxyServiceHandler(service.Server(), new(dbRpc.DBProxy))
	if err != nil {
		log.Fatalf("failed to register dbproxy service: %v", err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		log.Fatalf("failed to run dbproxy service: %v", err)
	}

}

func main() {
	startRpcService()
}
