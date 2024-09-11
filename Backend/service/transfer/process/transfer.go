package process

import (
	"bufio"
	cfg "cloud_distributed_storage/Backend/config"
	"cloud_distributed_storage/Backend/mq"
	dbcli "cloud_distributed_storage/Backend/service/dbproxy/client"
	"cloud_distributed_storage/Backend/store/s3"
	"encoding/json"
	"log"
	"os"
)

// Transfer : 处理文件转移
func Transfer(msg []byte) bool {
	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	s3Client := s3.GetS3Client()
	bucketBasics := s3.BucketBasics{S3Client: s3Client}
	err = bucketBasics.UploadFile(
		cfg.S3_BUCKET_NAME,
		pubData.DestLocation,
		bufio.NewReader(fin))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	resp, err := dbcli.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !resp.Suc {
		log.Println("更新数据库异常，请检查:" + pubData.FileHash)
		return false
	}
	return true
}
