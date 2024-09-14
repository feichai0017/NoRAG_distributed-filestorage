package api

import (
	"cloud_distributed_storage/Backend/common"
	cfg "cloud_distributed_storage/Backend/config"
	dblayer "cloud_distributed_storage/Backend/database"
	"cloud_distributed_storage/Backend/meta"
	"cloud_distributed_storage/Backend/mq"
	"cloud_distributed_storage/Backend/store/ceph"
	"cloud_distributed_storage/Backend/store/s3"
	"cloud_distributed_storage/Backend/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

// UploadHandler: handle file upload
func UploadHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewRespMsg(-1, "Failed to get file", nil))
		return
	}
	defer file.Close()

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewRespMsg(-2, "User not authenticated", nil))
		return
	}

	fileMeta := meta.FileMeta{
		FileName: header.Filename,
		Location: cfg.TempLocalRootDir + header.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewRespMsg(-3, "Failed to create file", nil))
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewRespMsg(-4, "Failed to save file", nil))
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	// 判断存储策略
	var storageType string
	if isImportantFile(fileMeta) {
		storageType = "ceph"
	} else {
		storageType = "s3"
	}

	switch storageType {
	case "ceph":
		// 保存到Ceph
		data, _ := ioutil.ReadAll(newFile)
		cephPath := "/ceph/" + fileMeta.FileSha1
		err = ceph.PutObject("userfile", cephPath, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.NewRespMsg(-5, "Failed to save to Ceph", nil))
			return
		}
		fileMeta.Location = cephPath
	case "s3":
		// 保存到S3
		s3Path := "/s3/" + fileMeta.FileSha1
		if cfg.AsyncTransferEnable {
			data := mq.TransferData{
				FileHash:      fileMeta.FileSha1,
				CurLocation:   fileMeta.Location,
				DestLocation:  s3Path,
				DestStoreType: common.StoreS3,
			}
			pubData, _ := json.Marshal(data)
			pubSuc := mq.Publish(cfg.TransExchangeName, cfg.TransS3RoutingKey, pubData)
			if !pubSuc {
				// TODO: 重试当前消息
			}
		} else {
			s3Client := s3.GetS3Client()
			err = s3Client.PutObject(cfg.S3BucketName, s3Path, newFile, "")
			if err != nil {
				c.JSON(http.StatusInternalServerError, util.NewRespMsg(-6, "Failed to save to S3", nil))
				return
			}
		}
		fileMeta.Location = s3Path
	}

	// 更新文件元数据
	_, err = meta.UpdateFileMetaDB(fileMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewRespMsg(-7, "Failed to update file meta", nil))
		return
	}

	// 更新用户文件表
	userFile := dblayer.UserFile{
		UserID:   userID.(int64),
		FileHash: fileMeta.FileSha1,
		FileName: fileMeta.FileName,
		FileSize: fileMeta.FileSize,
	}
	suc := dblayer.OnUserFileUploadFinished(userFile)
	if suc {
		c.JSON(http.StatusOK, util.NewRespMsg(0, "Upload successful", nil))
	} else {
		c.JSON(http.StatusInternalServerError, util.NewRespMsg(-8, "Failed to update user file info", nil))
	}
}

// 判断文件是否为重要文件
func isImportantFile(fileMeta meta.FileMeta) bool {
	// 可以根据文件名、大小、类型等条件判断
	matched, _ := regexp.MatchString("VI$", fileMeta.FileName)
	if matched {
		return true
	}
	return false
}
