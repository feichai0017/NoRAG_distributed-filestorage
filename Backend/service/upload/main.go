package main

import (
	"cloud_distributed_storage/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/success", handler.HTTPInterceptor(handler.UploadSucHandler))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/meta/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	//redis upload api
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uploadpart", handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/usr/local/Distributed_system/cloud_distributed_storage/Backend/static"))))

	// 配置 /service 路由来服务HTML文件
	http.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/service.html")
	})
	http.HandleFunc("/service/show", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/usr/local/Distributed_system/cloud_distributed_storage/Backend/static/view/fileDisplay.html")
	})
	http.HandleFunc("/user/info", handler.UserInfoHandler)
	http.HandleFunc("/user/login", handler.SignInHandler)
	http.HandleFunc("/", handler.SignInHandler)
	http.HandleFunc("/user/signup", handler.HTTPInterceptor(handler.SignupHandler))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}
