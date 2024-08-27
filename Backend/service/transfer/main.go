package main

import (
	"bufio"
	"cloud_distributed_storage/Backend/config"
	dblayer "cloud_distributed_storage/Backend/database"
	"cloud_distributed_storage/Backend/mq"
	"cloud_distributed_storage/Backend/store/s3"
	"encoding/json"
	"log"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	log.Printf("Received message: %s\n", string(msg)) // 打印接收到的消息

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	filed, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	s3Client := s3.GetS3Client()
	bucketBasics := s3.BucketBasics{S3Client: s3Client}
	err = bucketBasics.UploadFile(config.S3_BUCKET_NAME, pubData.FileHash, bufio.NewReader(filed))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	pubData.DestLocation = "s3://" + config.S3_BUCKET_NAME + "/" + pubData.FileHash

	suc := dblayer.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	if !suc {
		return false
	}

	return true

}

func main() {
	log.Println("start monitoring the queue")
	mq.StartConsumer(
		config.TransS3QueueName,
		"transfer_s3",
		ProcessTransfer)
}
