package handler

import (
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
	"time"
)

// UploadHandler: handle file upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//	return html
		file, err := os.Open("/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/upload.html")
		if err != nil {
			log.Printf("Failed to read file: %s, error: %s\n", "./static/view/upload.html", err.Error())
			io.WriteString(w, "internal server error")
			return
		}
		log.Printf("Successfully read file: %s\n", "./static/view/upload.html")
		defer file.Close()

		data, err := io.ReadAll(file)

		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
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
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/queryfile.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
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
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMetas := meta.GetLastFileMetas(limitCnt)
	fileMetas, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// DownloadHandler: download the file
// TODO: download file from ceph and fix the bug
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fsha1 := query.Get("filehash")
	// 获取文件元信息
	fm := meta.GetFileMeta(fsha1)
	if fm.Location == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}
	// download file from s3
	s3Client := s3.GetS3Client()
	bucketBasics := s3.BucketBasics{S3Client: s3Client}
	err := bucketBasics.DownloadFile(cfg.S3_BUCKET_NAME, fm.FileSha1, fm.FileName)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get file info"))
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

	_, err = io.Copy(w, f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
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
	query := r.URL.Query()
	fileSha1 := query.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

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
