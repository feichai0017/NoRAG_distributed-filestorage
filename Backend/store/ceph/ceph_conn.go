package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

// GetCephConn 获取ceph连接
func GetCephConn() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	auth := aws.Auth{
		AccessKey: "XPEJZ2P808Y9UMU29GDO",
		SecretKey: "52Ybt4XFLekFOWXNfdgdXqUMNowtttupkX64IY42",
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
