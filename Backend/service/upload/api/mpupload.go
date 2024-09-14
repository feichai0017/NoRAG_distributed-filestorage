package api

import (
	rPool "cloud_distributed_storage/Backend/cache/redis"
	"cloud_distributed_storage/Backend/config"
	cfg "cloud_distributed_storage/Backend/config"
	dbcli "cloud_distributed_storage/Backend/service/dbproxy/client"
	"cloud_distributed_storage/Backend/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func init() {
	os.MkdirAll(config.TempPartRootDir, 0744)
}

// InitialMultipartUploadHandler : 初始化分块上传
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

	// 2. 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 4. 将初始化信息写入到redis缓存
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 5. 将响应初始化数据返回到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": upInfo,
		})
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	//	username := c.Request.FormValue("username")
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 2. 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := config.TempPartRootDir + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": 0,
				"msg":  "Upload part failed",
				"data": nil,
			})
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4. 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})
}

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(c *gin.Context) {
	// 1. 解析请求参数
	upid := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize := c.Request.FormValue("filesize")
	filename := c.Request.FormValue("filename")

	// 2. 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "服务错误",
				"data": nil,
			})
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
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "分块不完整",
				"data": nil,
			})
		return
	}

	// 4. TODO：合并分块, 可以将ceph当临时存储，合并时将文件写入ceph;
	// 也可以不用在本地进行合并，转移的时候将分块append到ceph/s3即可
	srcPath := config.TempPartRootDir + upid + "/"
	destPath := cfg.TempLocalRootDir + filehash
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
	mergeRes, err := util.ExecLinuxShell(cmd)
	if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "合并失败",
				"data": nil,
			})
		return
	}
	log.Println(mergeRes)

	// 5. 更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)

	fmeta := dbcli.FileMeta{
		FileSha1: filehash,
		FileName: filename,
		FileSize: int64(fsize),
		Location: destPath,
	}
	_, ferr := dbcli.OnFileUploadFinished(fmeta)
	_, uferr := dbcli.OnUserFileUploadFinished(username, fmeta)
	if ferr != nil || uferr != nil {
		log.Println(err)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "数据更新失败",
				"data": nil,
			})
		return
	}

	// 6. 响应处理结果
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})
}

// CancelUploadHandler : 取消上传
func CancelUploadHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 删除文件分块
	os.RemoveAll(config.TempPartRootDir + uploadID)

	// 清除redis缓存
	rConn.Do("DEL", "MP_"+uploadID)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK", "data": nil})
}

// MultipartUploadStatusHandler : 查询分块上传的状态
func MultipartUploadStatusHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadid")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "查询失败", "data": nil})
		return
	}

	ret := make(map[string]interface{})
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		ret[k] = v
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK", "data": ret})
}

// MultiDownloadHandler : 断点续传下载
func MultiDownloadHandler(c *gin.Context) {
	filehash := c.Query("filehash")
	username := c.Query("username")

	// 检查用户对文件的访问权限
	permResult, err := dbcli.CheckPermission(username, filehash)
	if err != nil || !permResult.Suc {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "Permission denied", "data": nil})
		return
	}

	fmetaResult, err := dbcli.GetFileMeta(filehash)
	if err != nil || !fmetaResult.Suc {
		c.JSON(http.StatusNotFound, gin.H{"code": -1, "msg": "File not found", "data": nil})
		return
	}

	fmeta := dbcli.ToTableFile(fmetaResult.Data)

	f, err := os.Open(fmeta.FileAddr.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "File cannot be opened", "data": nil})
		return
	}
	defer f.Close()

	fileSize := fmeta.FileSize.Int64
	fileName := fmeta.FileName.String

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Accept-Ranges", "bytes")

	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		var start, end int64
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Invalid range header", "data": nil})
			return
		}

		if end == 0 {
			end = fileSize - 1
		}

		if start > end || start < 0 || end >= fileSize {
			c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
		c.Status(http.StatusPartialContent)

		_, err = f.Seek(start, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "Failed to seek file", "data": nil})
			return
		}

		_, err = io.CopyN(c.Writer, f, end-start+1)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "Failed to copy file content", "data": nil})
			return
		}
	} else {
		c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
		_, err = io.Copy(c.Writer, f)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "Failed to copy file content", "data": nil})
			return
		}
	}

	// 更新下载次数等信息
	_, err = dbcli.UpdateUserFileDownloadCount(username, filehash)
	if err != nil {
		// 仅记录日志，不影响下载
		log.Printf("Failed to update download count: %v", err)
	}
}
