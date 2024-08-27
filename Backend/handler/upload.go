package handler

import (
	"bytes"
	"cloud_distributed_storage/Backend/common"
	cfg "cloud_distributed_storage/Backend/config"
	dblayer "cloud_distributed_storage/Backend/database"
	"cloud_distributed_storage/Backend/meta"
	"cloud_distributed_storage/Backend/mq"
	"cloud_distributed_storage/Backend/store/ceph"
	"cloud_distributed_storage/Backend/store/s3"
	"cloud_distributed_storage/Backend/util"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	S3 "gopkg.in/amz.v1/s3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// FileRequest 定义请求的结构
type FileRequest struct {
	FileHash string `json:"filehash"`
	Limit    int    `json:"limit"`
}

func GetUsernameFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("username not found in context")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return "", errors.New("username is not a string")
	}

	return usernameStr, nil
}

// UploadHandler: handle file upload
func UploadHandler(c *gin.Context) {
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to get data", nil))
		return
	}
	username, err := GetUsernameFromContext(c)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "User not authenticated", nil))
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "/tmp/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 判断存储策略
	var storageType string
	if isImportantFile(fileMeta) {
		storageType = "ceph"
	} else {
		storageType = "s3"
	}

	newFile, err := os.Create(fileMeta.Location)

	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to create file", nil))
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to save data into file", nil))
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	newFile.Seek(0, 0)

	if storageType == "ceph" {
		// save file into ceph

		data, _ := ioutil.ReadAll(newFile)
		bucket := ceph.GetCephBucket("userfile")
		cephPath := "/ceph/" + fileMeta.FileSha1
		err = bucket.Put(cephPath, data, "octet-stream", S3.PublicRead)
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to save data into Ceph", nil))
			return
		}
		fileMeta.Location = cephPath
	} else if storageType == "s3" {
		// save file into s3
		//s3Client := s3.GetS3Client()
		//bucketBasics := s3.BucketBasics{S3Client: s3Client}
		//err = bucketBasics.UploadFile(cfg.S3_BUCKET_NAME, fileMeta.FileSha1, newFile)
		//if err != nil {
		//	fmt.Printf("Failed to upload file to S3, err:%s\n", err.Error())
		//	w.Write([]byte("upload failed."))
		//	return
		//}
		data := mq.TransferData{
			FileHash:      fileMeta.FileSha1,
			CurLocation:   fileMeta.Location,
			DestLocation:  cfg.S3_BUCKET_NAME,
			DestStoreType: common.StoreType(2),
		}
		pubData, _ := json.Marshal(data)
		log.Printf("start to publish message to transcode: %s\n", pubData)
		pubSuc := mq.Publish(cfg.TransExchangeName, cfg.TransS3RoutingKey, pubData)
		if !pubSuc {
			//	TODO: retry current message
		}
	}

	//meta.UpdateFileMeta(fileMeta)
	_ = meta.UpdateFileMetaDB(fileMeta)

	// 更新用户文件表记录
	suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
	if suc {
		c.Redirect(http.StatusFound, "/file/upload/success")
	} else {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Upload failed", nil))
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

// UploadSucHandler:upload success
func UploadSucHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/uploadSuccess.html")
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to get data", nil))
		return
	}
	c.Data(http.StatusOK, "text/html", data)
}

// GetFileMetaHandler: get meta info of file
func GetFileMetaHandler(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Invalid JSON payload", nil))
		return
	}

	fMeta, err := meta.GetFileMetaDB(req.FileHash)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to get file meta", nil))
		return
	}
	c.JSON(http.StatusOK, fMeta)
}

// FileQueryHandler: get meta info of file
func FileQueryHandler(c *gin.Context) {
	var req FileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Invalid JSON payload", nil))
		return
	}

	if req.Limit <= 0 {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Invalid limit", nil))
		return
	}

	username, err := GetUsernameFromContext(c)
	log.Printf("username: %s\n", username)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "User not authenticated", nil))
		return
	}

	fileMetas, err := dblayer.QueryUserFileMetas(username, req.Limit)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to query user files", nil))
		return
	}
	c.JSON(http.StatusOK, fileMetas)
}

// DownloadHandler: download the file
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Query("filehash")

	fm, err := meta.GetFileMetaDB(fsha1)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to get file meta", nil))
		return
	}

	var reader io.ReadCloser
	var fileSize int64

	switch {
	case strings.HasPrefix(fm.Location, "ceph/"):
		bucket := ceph.GetCephBucket("userfile")
		cephPath := strings.TrimPrefix(fm.Location, "ceph/")
		data, err := bucket.Get(cephPath)
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to get file from Ceph", nil))
			return
		}
		reader = ioutil.NopCloser(bytes.NewReader(data))
		fileSize = int64(len(data))

	case strings.HasPrefix(fm.Location, "s3/"):
		s3Client := s3.GetS3Client()
		bucketBasics := s3.BucketBasics{S3Client: s3Client}
		tempFile, err := ioutil.TempFile("", "s3-download-")
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to create temp file", nil))
			return
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		err = bucketBasics.DownloadFile(cfg.S3_BUCKET_NAME, fm.FileSha1, tempFile.Name())
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to download file from S3", nil))
			return
		}
		reader, err = os.Open(tempFile.Name())
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to open file", nil))
			return
		}
		fileSize = fm.FileSize

	default:
		reader, err = os.Open(fm.Location)
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to open file", nil))
			return
		}
		defer reader.Close()
		fi, err := reader.(*os.File).Stat()
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to get file info", nil))
			return
		}
		fileSize = fi.Size()
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fm.FileName+"\"")
	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))

	bufSize := 4 * 1024 * 1024
	buf := make([]byte, bufSize)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to read file", nil))
			return
		}
		if n == 0 {
			break
		}
		if _, err := c.Writer.Write(buf[:n]); err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to write response", nil))
			return
		}
		if f, ok := c.Writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// FileMetaUpdateHandler: update the filename of filemeta
func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.PostForm("op")
	fileSha1 := c.PostForm("filehash")
	newFileName := c.PostForm("filename")

	if opType != "0" {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Invalid operation type", nil))
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	c.JSON(http.StatusOK, curFileMeta)
}

// FileDeleteHandler: remove file from filemeta and local
func FileDeleteHandler(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Invalid JSON payload", nil))
		return
	}

	if len(req.FileHash) <= 0 {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Invalid file hash", nil))
		return
	}

	fMeta, err := meta.GetFileMetaDB(req.FileHash)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to get file meta", nil))
		return
	}

	switch {
	case strings.HasPrefix(fMeta.Location, "ceph/"):
		bucket := ceph.GetCephBucket("userfile")
		err = bucket.Del(fMeta.Location)

	case strings.HasPrefix(fMeta.Location, "s3/"):
		s3Client := s3.GetS3Client()
		bucketBasics := s3.BucketBasics{S3Client: s3Client}
		fileSha1List := []string{req.FileHash}
		err = bucketBasics.DeleteObjects(cfg.S3_BUCKET_NAME, fileSha1List)

	default:
		err = os.Remove(fMeta.Location)
	}

	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to remove file", nil))
		return
	}

	err = meta.RemoveFileMeta(req.FileHash)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to remove file meta", nil))
		return
	}

	c.JSON(http.StatusOK, util.NewRespMsg(0, "Delete success", nil))
}

func TryFastUploadHandler(c *gin.Context) {
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	filesize, _ := strconv.Atoi(c.PostForm("filesize"))

	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusOK, util.RespMsg{Code: -1, Msg: "failed"})
		return
	}

	if fileMeta.FileSha1 == "" {
		c.JSON(http.StatusOK, util.RespMsg{Code: -1, Msg: "failed"})
		return
	}

	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		c.JSON(http.StatusOK, util.RespMsg{Code: 0, Msg: "success upload immediately"})
	} else {
		c.JSON(http.StatusOK, util.RespMsg{Code: -2, Msg: "fail to upload immediately"})
	}
}
