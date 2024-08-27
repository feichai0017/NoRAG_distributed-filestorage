package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// 加载 .env 文件
func init() {
	err := godotenv.Load("/usr/local/Distributed_system/cloud_distributed_storage/Backend/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println(err.Error())
	}
}

var (
	// CephAccessKey : 访问Key
	CephAccessKey = os.Getenv("CEPH_ACCESS_KEY")
	// CephSecretKey : 访问密钥
	CephSecretKey = os.Getenv("CEPH_SECRET_KEY")
	// CephGWEndpoint : gateway地址
	CephGWEndpoint = "http://172.20.0.15:7480"
)
