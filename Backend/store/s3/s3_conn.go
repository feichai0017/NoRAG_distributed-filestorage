package s3

import (
	"context"
	"log"

	cfg "cloud_distributed_storage/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

// InitS3Client: Initialize the S3 client
func InitS3Client() {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.S3_REGION),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3_ACCESS_KEY, cfg.S3_SECRET_KEY, "")),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	s3Client = s3.NewFromConfig(awsConfig)
}

// GetS3Client: Get the S3 client
func GetS3Client() *s3.Client {
	if s3Client == nil {
		InitS3Client()
	}
	return s3Client
}
