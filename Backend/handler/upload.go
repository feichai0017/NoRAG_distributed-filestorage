package handler

import (
	"bytes"
	"cloud_distributed_storage/common"
	cfg "cloud_distributed_storage/config"
	dblayer "cloud_distributed_storage/database"
	"cloud_distributed_storage/meta"
	"cloud_distributed_storage/mq"
	"cloud_distributed_storage/store/ceph"
	"cloud_distributed_storage/store/s3"
	"cloud_distributed_storage/util"
	"encoding/json"
	"fmt"
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

// UploadHandler: handle file upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//	receive file stream
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err:%s\n", err.Error())
			return
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
			fmt.Printf("Failed to create file, err:%s\n", err.Error())
			return
		}

		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
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
				fmt.Printf("Failed to save data into ceph, err:%s\n", err.Error())
				w.Write([]byte("upload failed."))
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

		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/file/upload/success", http.StatusFound)

		}
		w.Write([]byte("upload failed."))
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
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/uploadSuccess.html")
	if err != nil {
		io.WriteString(w, "internal server error")
		return
	}
	io.WriteString(w, string(data))
}

// GetFileMetaHandler: get meta info of file
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 解析 JSON 请求
	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(req.FileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileQueryHandler: get meta info of file
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 解析 JSON 请求
	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// 验证 limit 值
	if req.Limit <= 0 {
		http.Error(w, `{"error": "Invalid limit value"}`, http.StatusBadRequest)
		return
	}

	// 从上下文中获取用户名
	username, ok := r.Context().Value("username").(string)
	if !ok {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}
	//fileMetas := meta.GetLastFileMetas(limitCnt)
	fileMetas, err := dblayer.QueryUserFileMetas(username, req.Limit)
	if err != nil {
		log.Printf("Error querying user file metas: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, err := json.Marshal(fileMetas)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// DownloadHandler: download the file
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fsha1 := query.Get("filehash")

	// 获取文件元信息
	fm, err := meta.GetFileMetaDB(fsha1)
	if err != nil {
		log.Printf("Failed to get file meta: %v", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}

	var reader io.ReadCloser
	var fileSize int64

	// 根据存储位置选择下载方式
	switch {
	case strings.HasPrefix(fm.Location, "ceph/"):
		// 从Ceph下载
		bucket := ceph.GetCephBucket("userfile")
		cephPath := strings.TrimPrefix(fm.Location, "ceph/")
		data, err := bucket.Get(cephPath)
		if err != nil {
			log.Printf("Failed to get file from Ceph: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		reader = ioutil.NopCloser(bytes.NewReader(data))
		fileSize = int64(len(data))

	case strings.HasPrefix(fm.Location, "s3/"):
		// 从S3下载
		s3Client := s3.GetS3Client()
		bucketBasics := s3.BucketBasics{S3Client: s3Client}
		tempFile, err := ioutil.TempFile("", "s3-download-")
		if err != nil {
			log.Printf("Failed to create temp file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		err = bucketBasics.DownloadFile(cfg.S3_BUCKET_NAME, fm.FileSha1, tempFile.Name())
		if err != nil {
			log.Printf("Failed to download file from S3: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		reader, err = os.Open(tempFile.Name())
		if err != nil {
			log.Printf("Failed to open temp file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fileSize = fm.FileSize

	default:
		// 从本地文件系统下载
		reader, err = os.Open(fm.Location)
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer reader.Close()
		fi, err := reader.(*os.File).Stat()
		if err != nil {
			log.Printf("Failed to get file info: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fileSize = fi.Size()
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))

	// 使用缓冲写入响应
	bufSize := 4 * 1024 * 1024 // 4MB buffer
	buf := make([]byte, bufSize)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Error reading file: %v", err)
			return
		}
		if n == 0 {
			break
		}
		if _, err := w.Write(buf[:n]); err != nil {
			log.Printf("Error writing to response: %v", err)
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// FileMetaUpdateHandler: update the filename of filemeta
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler: remove file from filemeta and local
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 解析 JSON 请求
	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	if len(req.FileHash) <= 0 {
		http.Error(w, `{"error": "Invalid limit value"}`, http.StatusBadRequest)
		return
	}
	log.Printf("file hash: %s", req.FileHash)

	fMeta, err := meta.GetFileMetaDB(req.FileHash)
	if err != nil {
		log.Printf("Failed to get file meta: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 根据存储位置选择删除方式
	switch {
	case strings.HasPrefix(fMeta.Location, "ceph/"):
		// 从Ceph删除
		bucket := ceph.GetCephBucket("userfile")
		err = bucket.Del(fMeta.Location)

	case strings.HasPrefix(fMeta.Location, "s3/"):
		// 从S3删除
		s3Client := s3.GetS3Client()
		bucketBasics := s3.BucketBasics{S3Client: s3Client}
		// 将fileSha1转换为[]string
		fileSha1List := []string{req.FileHash}
		err = bucketBasics.DeleteObjects(cfg.S3_BUCKET_NAME, fileSha1List)

	default:
		// 从本地文件系统删除
		err = os.Remove(fMeta.Location)
	}

	if err != nil {
		log.Printf("Failed to delete file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 从数据库中删除文件元信息
	err = meta.RemoveFileMeta(req.FileHash)
	if err != nil {
		log.Printf("Failed to remove file meta from DB: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fileMeta.FileSha1 == "" {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "failed",
		}
		w.Write(resp.JSONBytes())
		return
	}
	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "success upload immediately",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "fail to upload immediately",
	}
	w.Write(resp.JSONBytes())
	return
}
