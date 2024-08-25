package ceph

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
	"log"
	"os"
)

var cephConn *s3.S3

// GetCephConn 获取ceph连接
func GetCephConn() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}
	// 加载 .env 文件
	err := godotenv.Load("/usr/local/Distributed_system/cloud_distributed_storage/Backend/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println(err.Error())
	}

	auth := aws.Auth{
		AccessKey: os.Getenv("CEPH_ACCESS_KEY"),
		SecretKey: os.Getenv("CEPH_SECRET_KEY"),
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "http://172.20.0.15:7480",
		S3Endpoint:           "http://172.20.0.15:7480",
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}

	return s3.New(auth, curRegion)
}

// GetCephBucket 获取ceph bucket
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConn()
	return conn.Bucket(bucket)
}
