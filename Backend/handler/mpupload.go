package handler

import (
	rPool "cloud_distributed_storage/Backend/cache/redis"
	"cloud_distributed_storage/Backend/config"
	dblayer "cloud_distributed_storage/Backend/database"
	"cloud_distributed_storage/Backend/meta"
	"cloud_distributed_storage/Backend/util"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// 初始化分块上传
type MultiPartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func init() {
	os.MkdirAll(config.TempPartRootDir, 0744)
}

// InitialMultipartUploadHandler: initialize multipart upload
func InitialMultipartUploadHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "params invalid",
			})
		return
	}

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	upInfo := MultiPartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", upInfo))
}

// UploadPartHandler:上传分块
func UploadPartHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")
	chunkIndex := c.PostForm("index")
	expectedHash := c.PostForm("filehash")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Upload part failed", nil))
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	hash := sha1.New()
	for {
		n, err := c.Request.Body.Read(buf)
		if n > 0 {
			hash.Write(buf[:n])
			fd.Write(buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Upload part failed", nil))
			return
		}
	}

	calculatedHash := hex.EncodeToString(hash.Sum(nil))
	if calculatedHash != expectedHash {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Hash mismatch", nil))
		return
	}

	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", nil))
}

// CompleteUploadHandler:通知上传合并
func CompleteUploadHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filesize := c.PostForm("filesize")
	filename := c.PostForm("filename")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Complete upload failed", nil))
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))

		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount += 1
		}
	}
	if totalCount != chunkCount {
		c.JSON(http.StatusOK, util.NewRespMsg(-2, "Invalid request", nil))
		return
	}

	mergeFile, err := os.Create("/data/" + filehash)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Complete upload failed", nil))
		return
	}
	defer mergeFile.Close()

	for i := 0; i < totalCount; i++ {
		fpath := "/data/" + uploadID + "/" + strconv.Itoa(i)
		chunkFile, err := os.Open(fpath)
		if err != nil {
			c.JSON(http.StatusOK, util.NewRespMsg(-1, "Complete upload failed", nil))
			return
		}
		defer chunkFile.Close()
		buf := make([]byte, 1024*1024)
		for {
			n, err := chunkFile.Read(buf)
			if n > 0 {
				mergeFile.Write(buf[:n])
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				c.JSON(http.StatusOK, util.NewRespMsg(-1, "Complete upload failed", nil))
				return
			}
		}
	}

	_, err = rConn.Do("DEL", "MP_"+uploadID)
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Failed to clear upload info from Redis", nil))
		return
	}

	fileSizeInt, _ := strconv.ParseInt(filesize, 10, 64)
	fileMeta := meta.FileMeta{
		FileSha1: filehash,
		FileName: filename,
		FileSize: fileSizeInt,
		Location: "/data/" + filehash,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	_ = meta.UpdateFileMetaDB(fileMeta)
	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, fileSizeInt)

	if suc {
		c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", nil))
	} else {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Upload failed", nil))
	}
}

// CancelUploadHandler:取消上传
func CancelUploadHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	os.RemoveAll("/data/" + uploadID)
	rConn.Do("DEL", "MP_"+uploadID)

	c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", nil))
}

// MultipartUploadStatusHandler:查询分块上传的状态
func MultipartUploadStatusHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		c.JSON(http.StatusOK, util.NewRespMsg(-1, "Multipart upload failed", nil))
		return
	}
	ret := make(map[string]interface{})
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		ret[k] = v
	}
	c.JSON(http.StatusOK, util.NewRespMsg(0, "OK", ret))
}

// TODO: 断点续传
func MultiDownloadHandler(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	//fsha1 := r.Form.Get("filehash")
	//username := r.Form.Get("username")
	//
	//fm, err := meta.GetFileMetaDB(fsha1)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//f, err := os.Open(fm.Location)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//defer f.Close()
	//
	//data, err := ioutil.ReadAll(f)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "application/octet-stream")
	//w.Header().Set("Content-Disposition", "attachment; filename="+fm.FileName)
	//w.Header().Set("Content-Length", strconv.FormatInt(fm.FileSize, 10))
	//w.Header().Set("Accept-Ranges", "bytes")
	//
	//rangeHeader := r.Header.Get("Range")
	//if rangeHeader != "" {
	//	// 处理断点续传
	//	var start, end int64
	//	fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	//	if end == 0 {
	//		end = fm.FileSize - 1
	//	}
	//	if start > end || start < 0 || end >= fm.FileSize {
	//		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
	//		return
	//	}
	//	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fm.FileSize))
	//	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	//	w.WriteHeader(http.StatusPartialContent)
	//	w.Write(data[start : end+1])
	//} else {
	//	// 完整下载
	//	w.Write(data)
	//}

	// 更新下载次数等信息
	//r.ParseForm()
	//username := r.Form.Get("username")
	//dblayer.UpdateUserFileDownloadCount(username, fsha1)
}
