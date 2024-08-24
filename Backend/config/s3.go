package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	// Load .env file
	err := godotenv.Load("cloud_distributed_storage/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var (
	// S3_BUCKET_NAME: returns the AWS S3 Bucket Name from environment variables
	S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME")

	// S3_REGION: returns the AWS S3 Region from environment variables
	S3_REGION = os.Getenv("S3_REGION")

	// S3_ENDPOINT: returns the AWS S3 Endpoint from environment variables
	S3_ENDPOINT = os.Getenv("S3_ENDPOINT")

	// S3_ACCESS_KEY returns the AWS Access Key from environment variables
	S3_ACCESS_KEY = os.Getenv("S3_ACCESS_KEY")

	// S3_SECRET_KEY returns the AWS Secret Access Key from environment variables
	S3_SECRET_KEY = os.Getenv("S3_SECRET_KEY")
)
