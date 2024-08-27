package ceph

import (
	"cloud_distributed_storage/Backend/config"
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
		AccessKey: config.CephAccessKey,
		SecretKey: config.CephSecretKey,
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          config.CephGWEndpoint,
		S3Endpoint:           config.CephGWEndpoint,
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
