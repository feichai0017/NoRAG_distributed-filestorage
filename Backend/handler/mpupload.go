package handler

import (
	rPool "cloud_distributed_storage/cache/redis"
	dblayer "cloud_distributed_storage/database"
	"cloud_distributed_storage/util"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
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

func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}
	//connect redis pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	upInfo := MultiPartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, //5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}
	// 将初始化信息写入redis
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// return response
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())

}

// UploadPartHandler:上传分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")
	expectedHash := r.Form.Get("filehash") // 获取预期的哈希值
	// get a connection from redis pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 获取文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	// 读取分块内容并计算sha1
	buf := make([]byte, 1024*1024)
	hash := sha1.New()
	for {
		n, err := r.Body.Read(buf)
		if n > 0 {
			hash.Write(buf[:n])
			fd.Write(buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
			return
		}
	}

	calculatedHash := hex.EncodeToString(hash.Sum(nil))
	if calculatedHash != expectedHash {
		w.Write(util.NewRespMsg(-1, "Hash mismatch", nil).JSONBytes())
		return
	}

	// 更新redis状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler:通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	// get a connection from redis pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 获取分块信息
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil || len(data) == 0 {
		w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))

		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chikidx_") && v == "1" {
			chunkCount += 1
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "Invalid request", nil).JSONBytes())
		return
	}

	// 合并分块
	mergeFile, err := os.Create("/data/" + filehash)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
		return
	}
	defer mergeFile.Close()

	for i := 0; i < totalCount; i++ {
		fpath := "/data/" + uploadID + "/" + strconv.Itoa(i)
		chunkFile, err := os.Open(fpath)
		if err != nil {
			w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
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
				w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
				return
			}
		}
	}
	// 清除Redis中关于这个文件的分块信息
	_, err = rConn.Do("DEL", "MP_"+uploadID)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Failed to clear upload info from Redis", nil).JSONBytes())
		return
	}

	// upload to database
	fsiz, err := strconv.Atoi(filesize)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsiz), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsiz))

	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CancelUploadHandler:取消上传
func CancelUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")

	// get a connection from redis pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 删除分块文件
	os.RemoveAll("/data/" + uploadID)
	// 删除redis记录
	rConn.Do("DEL", "MP_"+uploadID)

	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// MultipartUploadStatusHandler:查询分块上传的状态
func MultipartUploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")

	// get a connection from redis pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Multipart upload failed", nil).JSONBytes())
		return
	}
	ret := make(map[string]interface{})
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		ret[k] = v
	}
	w.Write(util.NewRespMsg(0, "OK", ret).JSONBytes())
}

//TODO: 断点下载
